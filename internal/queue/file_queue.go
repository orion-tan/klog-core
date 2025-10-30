package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"klog-backend/internal/cache"
	"klog-backend/internal/utils"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis Stream名称
	fileDeleteStream = "file:delete:stream"
	// 消费者组名称
	consumerGroup = "file-delete-group"
	// 消费者名称
	consumerName = "file-delete-consumer"
	// 最大重试次数
	maxRetries = 3
)

// DeleteTask 文件删除任务
type DeleteTask struct {
	FilePath   string    `json:"file_path"`
	Timestamp  time.Time `json:"timestamp"`
	RetryCount int       `json:"retry_count"`
}

// FileQueue 文件删除队列
type FileQueue struct {
	useRedis bool
}

// NewFileQueue 创建文件删除队列
func NewFileQueue() *FileQueue {
	return &FileQueue{
		useRedis: cache.RedisClient != nil,
	}
}

// PublishDeleteTask 发布文件删除任务
// @filePath 文件路径
// @return 错误
func (q *FileQueue) PublishDeleteTask(filePath string) error {
	task := DeleteTask{
		FilePath:   filePath,
		Timestamp:  time.Now(),
		RetryCount: 0,
	}

	if q.useRedis {
		// 使用Redis Streams发布消息
		return q.publishToRedis(task)
	}

	// 降级处理：直接在goroutine中异步删除
	go q.deleteFileWithRetry(filePath, maxRetries)
	return nil
}

// publishToRedis 将删除任务发布到Redis Streams
// @task 删除任务
// @return 错误
func (q *FileQueue) publishToRedis(task DeleteTask) error {
	taskData, err := json.Marshal(task)
	if err != nil {
		utils.SugarLogger.Errorf("序列化删除任务失败: %v", err)
		return err
	}

	ctx := context.Background()
	_, err = cache.RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: fileDeleteStream,
		Values: map[string]interface{}{
			"task": string(taskData),
		},
	}).Result()

	if err != nil {
		utils.SugarLogger.Errorf("发布删除任务到Redis失败: %v, 降级为直接异步删除", err)
		// 发布失败时降级处理
		go q.deleteFileWithRetry(task.FilePath, maxRetries)
		return err
	}

	utils.SugarLogger.Infof("文件删除任务已发布到队列: %s", task.FilePath)
	return nil
}

// StartConsumer 启动消费者
// @ctx 上下文，用于控制消费者生命周期
func (q *FileQueue) StartConsumer(ctx context.Context) {
	if !q.useRedis {
		utils.SugarLogger.Info("Redis未配置，文件删除队列消费者未启动（降级模式）")
		return
	}

	// 确保消费者组存在
	if err := q.ensureConsumerGroup(); err != nil {
		utils.SugarLogger.Errorf("创建消费者组失败: %v", err)
		return
	}

	utils.SugarLogger.Info("文件删除队列消费者已启动")

	// 启动消费循环
	go q.consumeLoop(ctx)
}

// ensureConsumerGroup 确保消费者组存在
// @return 错误
func (q *FileQueue) ensureConsumerGroup() error {
	ctx := context.Background()

	// 尝试创建消费者组
	err := cache.RedisClient.XGroupCreateMkStream(ctx, fileDeleteStream, consumerGroup, "0").Err()
	if err != nil {
		// 如果组已存在，忽略错误
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			return nil
		}
		return err
	}

	return nil
}

// consumeLoop 消费循环
// @ctx 上下文
func (q *FileQueue) consumeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			utils.SugarLogger.Info("文件删除队列消费者停止")
			return
		default:
			// 从流中读取消息
			streams, err := cache.RedisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: consumerName,
				Streams:  []string{fileDeleteStream, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					// 没有新消息，继续循环
					continue
				}
				utils.SugarLogger.Errorf("读取消息失败: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// 处理消息
			for _, stream := range streams {
				for _, message := range stream.Messages {
					q.processMessage(ctx, message)
				}
			}
		}
	}
}

// processMessage 处理单条消息
// @ctx 上下文
// @message 消息
func (q *FileQueue) processMessage(ctx context.Context, message redis.XMessage) {
	taskJSON, ok := message.Values["task"].(string)
	if !ok {
		utils.SugarLogger.Errorf("消息格式错误: %v", message.ID)
		q.ackMessage(ctx, message.ID)
		return
	}

	var task DeleteTask
	if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
		utils.SugarLogger.Errorf("反序列化任务失败: %v", err)
		q.ackMessage(ctx, message.ID)
		return
	}

	// 执行文件删除
	if err := q.deleteFile(task.FilePath); err != nil {
		// 删除失败，检查重试次数
		if task.RetryCount < maxRetries {
			task.RetryCount++
			utils.SugarLogger.Warnf("删除文件失败 (重试 %d/%d): %s, 错误: %v",
				task.RetryCount, maxRetries, task.FilePath, err)

			// 重新发布任务
			if err := q.publishToRedis(task); err != nil {
				utils.SugarLogger.Errorf("重新发布任务失败: %v", err)
			}
		} else {
			utils.SugarLogger.Errorf("删除文件失败，已达最大重试次数: %s, 错误: %v",
				task.FilePath, err)
		}
	} else {
		utils.SugarLogger.Infof("成功删除文件: %s", task.FilePath)
	}

	// 确认消息已处理
	q.ackMessage(ctx, message.ID)
}

// ackMessage 确认消息
// @ctx 上下文
// @messageID 消息ID
func (q *FileQueue) ackMessage(ctx context.Context, messageID string) {
	if err := cache.RedisClient.XAck(ctx, fileDeleteStream, consumerGroup, messageID).Err(); err != nil {
		utils.SugarLogger.Errorf("确认消息失败: %v", err)
	}
}

// deleteFile 删除文件
// @filePath 文件路径
// @return 错误
func (q *FileQueue) deleteFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，视为成功
		return nil
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// deleteFileWithRetry 带重试的文件删除（用于降级模式）
// @filePath 文件路径
// @maxRetries 最大重试次数
func (q *FileQueue) deleteFileWithRetry(filePath string, maxRetries int) {
	for i := 0; i <= maxRetries; i++ {
		if err := q.deleteFile(filePath); err != nil {
			if i < maxRetries {
				utils.SugarLogger.Warnf("删除文件失败 (重试 %d/%d): %s, 错误: %v",
					i+1, maxRetries, filePath, err)
				time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
				continue
			}
			utils.SugarLogger.Errorf("删除文件失败，已达最大重试次数: %s, 错误: %v",
				filePath, err)
			return
		}
		utils.SugarLogger.Infof("成功删除文件: %s", filePath)
		return
	}
}
