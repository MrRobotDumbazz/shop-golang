package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

type Redis interface {
	InitRedis() *redis.Client
	SetToken(ctx context.Context, SID int, token string)
	GetToken(ctx context.Context, ID int) (string, error)
	DeleteToken(ctx context.Context, ID int)
	ExpireToken(ctx context.Context, ID int)
}

type RedisRepository struct {
	Redis
}

func (r *RedisClient) InitRedis() *redis.Client {
	addr := fmt.Sprintf("localhost:6379")
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func newRedisRepository(conn *redis.Client) *RedisClient {
	return &RedisClient{
		client: conn,
	}
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		Redis: newRedisRepository(client),
	}
}

func (r *RedisClient) SetToken(ctx context.Context, SID int, token string) {
	err := r.client.Set(ctx, fmt.Sprintf("token-%d", SID), token, time.Minute*10)
	if err != nil {
		log.Printf("error creating token in redis %v", err)
		return
	}
}

func (r *RedisClient) GetToken(ctx context.Context, ID int) (string, error) {
	token, err := r.client.Get(ctx, fmt.Sprint("token-%d", ID)).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *RedisClient) DeleteToken(ctx context.Context, ID int) {
	err := r.client.Del(ctx, fmt.Sprintf("token-%d", ID))
	if err != nil {
		log.Printf("error deleting token in redis %v", err)
		return
	}
}

func (r *RedisClient) ExpireToken(ctx context.Context, ID int) {
	err := r.client.Expire(ctx, fmt.Sprintf("token-%d", ID), 10)
	if err != nil {
		log.Printf("error expiring token in redis %v", err)
		return
	}
}

// func (redisCli *RedisCli) SetValue(key string, value string, expiration ...interface{}) error {
// 	_, err := redisCli.conn.Do("SET", key, value)

// 	if err == nil && expiration != nil {
// 		redisCli.conn.Do("EXPIRE", key, expiration[0])
// 	}

// 	return err
// }

// func (redisCli *RedisCli) GetValue(key string) (interface{}, error) {
// 	return redisCli.conn.Do("GET", key)
// }
