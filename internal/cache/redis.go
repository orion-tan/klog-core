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

// GetOrSet 获取缓存，如果不存在则执行loader函数并缓存结果
// 实现Cache-Aside模式，防止缓存穿透
// @key 缓存键
// @dest 目标对象（用于接收缓存数据）
// @expiration 过期时间
// @loader 数据加载函数（当缓存不存在时调用）
// @return 错误
func GetOrSet(key string, dest interface{}, expiration time.Duration, loader func() (interface{}, error)) error {
	if RedisClient == nil {
		// Redis未启用时直接调用loader
		data, err := loader()
		if err != nil {
			return err
		}
		// 将loader返回的数据赋值给dest（通过JSON序列化/反序列化）
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonData, dest)
	}

	// 尝试从缓存获取
	err := Get(key, dest)
	if err == nil {
		return nil // 缓存命中
	}

	if err != redis.Nil {
		return err // 发生错误
	}

	// 缓存未命中，调用loader加载数据
	data, err := loader()
	if err != nil {
		return err
	}

	// 设置缓存
	if err := Set(key, data, expiration); err != nil {
		return err
	}

	// 将数据赋值给dest
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, dest)
}

// SetMany 批量设置缓存
// @items map[key]value 键值对
// @expiration 过期时间
// @return 错误
func SetMany(items map[string]interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil // Redis未启用时跳过
	}

	pipe := RedisClient.Pipeline()
	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		pipe.Set(ctx, key, data, expiration)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetMany 批量获取缓存
// @keys 键列表
// @return map[key]value（仅包含存在的键）, 错误
func GetMany(keys []string) (map[string]string, error) {
	if RedisClient == nil {
		return make(map[string]string), nil // Redis未启用时返回空map
	}

	pipe := RedisClient.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[string]string)
	for i, cmd := range cmds {
		val, err := cmd.Result()
		if err == nil {
			result[keys[i]] = val
		}
	}

	return result, nil
}

// Increment 原子递增
// @key 键
// @delta 增量（可以为负数实现递减）
// @return 递增后的值, 错误
func Increment(key string, delta int64) (int64, error) {
	if RedisClient == nil {
		return 0, nil // Redis未启用时返回0
	}

	return RedisClient.IncrBy(ctx, key, delta).Result()
}

// Decrement 原子递减
// @key 键
// @delta 减量
// @return 递减后的值, 错误
func Decrement(key string, delta int64) (int64, error) {
	return Increment(key, -delta)
}

// SetNX 仅当key不存在时设置（分布式锁实现基础）
// @key 键
// @value 值
// @expiration 过期时间
// @return 是否设置成功, 错误
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if RedisClient == nil {
		return false, nil // Redis未启用时返回false
	}

	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return RedisClient.SetNX(ctx, key, data, expiration).Result()
}

// GetTTL 获取key的剩余生存时间
// @key 键
// @return 剩余TTL（-2表示不存在，-1表示无过期时间）, 错误
func GetTTL(key string) (time.Duration, error) {
	if RedisClient == nil {
		return 0, nil // Redis未启用时返回0
	}

	return RedisClient.TTL(ctx, key).Result()
}

// Expire 设置key的过期时间
// @key 键
// @expiration 过期时间
// @return 是否设置成功, 错误
func Expire(key string, expiration time.Duration) (bool, error) {
	if RedisClient == nil {
		return false, nil // Redis未启用时返回false
	}

	return RedisClient.Expire(ctx, key, expiration).Result()
}

// DeleteMany 批量删除缓存
// @keys 键列表
// @return 错误
func DeleteMany(keys []string) error {
	if RedisClient == nil {
		return nil // Redis未启用时跳过
	}

	if len(keys) == 0 {
		return nil
	}

	return RedisClient.Del(ctx, keys...).Err()
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
