package database

import (
	"context"
	"github.com/redis/rueidis"
	"log"
	"showmaster/src/config"
)

// GetRedisDatabase returns a rueidis.Client with the cfg config.Config provided
func GetRedisDatabase(cfg *config.Config) *rueidis.Client {
	// Get a new client
	client, err := rueidis.NewClient(rueidis.ClientOption{
		Username:    cfg.Database.Redis.User,
		Password:    cfg.Database.Redis.Password,
		InitAddress: []string{cfg.Database.Redis.Host + ":" + cfg.Database.Redis.Port},
		SelectDB:    cfg.Database.Redis.DB,
	})
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v\n", err)
		return nil
	}

	// Quick ping test
	ctx := context.Background()
	err = client.Do(ctx, client.B().Ping().Build()).Error()
	if err != nil {
		log.Fatalf("Error pinging Redis: %v\n", err)
		return nil
	}

	return &client
}
