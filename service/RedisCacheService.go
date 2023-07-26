package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCacheService struct {
	client  *redis.Client
	context context.Context
}

func NewRedisCacheService(redisHost string, redisPort int) *RedisCacheService {
	return &RedisCacheService{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
		}),
		context: context.Background(),
	}
}

func (c *RedisCacheService) Ping() error {
	// Verificar la conexión a Redis
	pong, err := c.client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	fmt.Println("Conexión exitosa a Redis:", pong)
	return nil
}

func (c *RedisCacheService) Set(key string, value interface{}) {
	err := c.client.Set(context.Background(), key, value, time.Hour*24).Err()
	if err != nil {
		fmt.Println("Error al guardar en caché:", err)
		return
	}
}

func (c *RedisCacheService) Get(key string) interface{} {
	value, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Println("Error al obtener el valor del caché: "+key, err)
		}
		return nil
	}

	return value
}
