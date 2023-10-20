package storages

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
)

func NewRedis(log *zap.Logger, ctx context.Context) (*redis.Client, error) {
	log.Info("Connecting to Redis...",
		zap.String("host", viper.GetString(config.RedisHost)),
		zap.String("port", viper.GetString(config.RedisPort)),
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(config.RedisHost) + ":" + viper.GetString(config.RedisPort),
		Password: viper.GetString(config.RedisPassword),
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Failed to create Redis connection, ", zap.Error(err))
		return nil, err
	}

	log.Info("Redis connection created successfully")
	return rdb, nil
}
