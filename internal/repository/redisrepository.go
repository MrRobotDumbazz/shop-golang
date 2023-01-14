package repository

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type RedisClient struct {
	client *redis.Client
}

type Redis interface {
	InitRedis() *redis.Client
	SetToken(SID int, token string)
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

func (r *RedisClient) SetToken(SID int, token string) {
	r.client.Set(fmt.Sprintf("token-%d", SID), token, 0)
}

func (r *RedisClient) GetToken(ID int, token string) (error, string) {
	token, err := r.client.Get(fmt.Sprint("token-%d", ID)).Result()
	if err != nil {
		return err, ""
	}
	return nil, token
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
