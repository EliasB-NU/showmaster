package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"showmaster/src/config"
)

func GetRedisDatabase(cfg *config.Config) *redis.Client {
	ctx := context.Background()

	rdp := redis.NewClient(&redis.Options{
		Addr:     cfg.Database.Redis.Host + ":" + cfg.Database.Redis.Port,
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	})

	_, err := rdp.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis database: %v", err)
		return nil
	}

	log.Println("Connected to Redis database")
	return rdp
}
