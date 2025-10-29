package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"klog-backend/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() error {
	if config.Cfg.Redis.Addr == "" {
		// Redis配置为空时跳过初始化
		return nil
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Password,
		DB:       0,
	})

	// 测试连接
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	return nil
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil // Redis未启用时跳过
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func Get(key string, dest interface{}) error {
	if RedisClient == nil {
		return redis.Nil // Redis未启用时返回nil
	}

	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func Delete(key string) error {
	if RedisClient == nil {
		return nil // Redis未启用时跳过
	}

	return RedisClient.Del(ctx, key).Err()
}

// DeleteByPattern 根据模式删除缓存
func DeleteByPattern(pattern string) error {
	if RedisClient == nil {
		return nil // Redis未启用时跳过
	}

	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := RedisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Exists 检查key是否存在
func Exists(key string) (bool, error) {
	if RedisClient == nil {
		return false, nil // Redis未启用时返回false
	}

	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// AddToBlacklist 将JWT token添加到黑名单（用于登出功能）
// @token JWT token字符串
// @expiration 过期时间（应与token的过期时间一致）
// @return 错误
func AddToBlacklist(token string, expiration time.Duration) error {
	if RedisClient == nil {
		return nil // Redis未启用时跳过
	}

	key := "jwt:blacklist:" + token
	return RedisClient.Set(ctx, key, "1", expiration).Err()
}

// IsInBlacklist 检查JWT token是否在黑名单中
// @token JWT token字符串
// @return 是否在黑名单中, 错误
func IsInBlacklist(token string) (bool, error) {
	if RedisClient == nil {
		return false, nil // Redis未启用时返回false
	}

	key := "jwt:blacklist:" + token
	exists, err := RedisClient.Exists(ctx, key).Result()
	return exists > 0, err
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
