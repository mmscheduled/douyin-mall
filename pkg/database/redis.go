package database

import (
    "context"
    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() error {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis 地址
        Password: "",               // 密码
        DB:       0,                // 数据库
    })

    _, err := RedisClient.Ping(context.Background()).Result()
    return err
}