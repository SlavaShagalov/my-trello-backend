package config

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	"github.com/spf13/viper"
)

// Postgres

func SetDefaultPostgresConfig() {
	viper.SetDefault(PostgresHost, "data-storage")
	viper.SetDefault(PostgresPort, 5432)
	viper.SetDefault(PostgresDB, "trello_db")
	viper.SetDefault(PostgresUser, "slava")
	viper.SetDefault(PostgresPassword, "1234")
	viper.SetDefault(PostgresSSLMode, "disable")
}

func SetTestPostgresConfig() {
	viper.SetDefault(PostgresHost, "data-storage")
	viper.SetDefault(PostgresPort, 5432)
	viper.SetDefault(PostgresDB, "trello_db")
	viper.SetDefault(PostgresUser, "slava")
	viper.SetDefault(PostgresPassword, "1234")
	viper.SetDefault(PostgresSSLMode, "disable")
}

// Redis

func SetDefaultRedisConfig() {
	viper.SetDefault(RedisHost, "sessions-storage")
	viper.SetDefault(RedisPort, "6379")
	viper.SetDefault(RedisPassword, "1234")
}

func SetTestRedisConfig() {
	viper.SetDefault(RedisHost, "sessions-storage")
	viper.SetDefault(RedisPort, "6379")
	viper.SetDefault(RedisPassword, "1234")
}

// S3

func SetDefaultS3Config() {
	viper.SetDefault(S3BucketName, "trello")
	viper.SetDefault(S3DefaultRegion, "ru-msk")
	viper.SetDefault(S3Endpoint, "http://hb.vkcs.cloud")
}

// Validation

func SetDefaultValidationConfig() {
	viper.SetDefault(MinUsernameLen, constants.MinUsernameLen)
	viper.SetDefault(MaxUsernameLen, constants.MaxUsernameLen)
}
