package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
)

func NewPostgres(log *zap.Logger) (*sql.DB, error) {
	log.Info("Connecting to db...",
		zap.String("host", viper.GetString(config.PostgresHost)),
		zap.Int("port", viper.GetInt(config.PostgresPort)),
		zap.String("dbname", viper.GetString(config.PostgresDB)),
	)

	params := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString(config.PostgresHost),
		viper.GetInt(config.PostgresPort),
		viper.GetString(config.PostgresUser),
		viper.GetString(config.PostgresDB),
		viper.GetString(config.PostgresPassword),
		viper.GetString(config.PostgresSSLMode),
	)

	db, err := sql.Open("postgres", params)
	if err != nil {
		log.Error("Failed to create DB connection", zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Error("Failed to connect to DB", zap.Error(err))
		return nil, err
	}

	log.Info("DB connection created successfully")
	return db, nil
}
