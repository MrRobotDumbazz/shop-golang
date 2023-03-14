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
	SetToken(ctx context.Context, SID int, token string) error
	GetToken(ctx context.Context, ID int) (string, error)
	GetTokenUUID(ctx context.Context) (int, error)
	DeleteToken(ctx context.Context, ID int)
	ExpireToken(ctx context.Context, ID int) error
}

type RedisRepository struct {
	Redis
}

func InitRedis() *redis.Client {
	addr := fmt.Sprintf("redis:6379")
	r := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	status := r.Ping(context.Background())
	err := status.Err()
	if err != nil {
		log.Fatal(err)
	}

	return r
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

func (r *RedisClient) SetToken(ctx context.Context, SID int, token string) error {
	log.Printf("Token UUID: %s", token)
	err := r.client.Set(ctx, fmt.Sprintf("token-%d", SID), token, time.Minute*10)
	// log.Println(err.Err().Error)
	log.Println(err)
	return nil
}

func (r *RedisClient) GetToken(ctx context.Context, ID int) (string, error) {
	log.Printf("Claims ID: %d", ID)
	token, _ := r.client.Get(ctx, fmt.Sprintf("token-%d", ID)).Result()
	return token, nil
}

func (r *RedisClient) DeleteToken(ctx context.Context, ID int) {
	err := r.client.Del(ctx, fmt.Sprintf("token-%d", ID))
	log.Println(err)
}

func (r *RedisClient) ExpireToken(ctx context.Context, ID int) error {
	err := r.client.Expire(ctx, fmt.Sprintf("token-%d", ID), 10*time.Minute)
	log.Println(err)
	return nil
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
