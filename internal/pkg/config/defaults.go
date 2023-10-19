package config

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	"github.com/spf13/viper"
)

// Postgres

func SetDefaultPostgresConfig() {
	viper.SetDefault(PostgresHost, "localhost")
	viper.SetDefault(PostgresPort, 5432)
	viper.SetDefault(PostgresDB, "judi_db")
	viper.SetDefault(PostgresUser, "judi")
	viper.SetDefault(PostgresPassword, "judi_pswd")
	viper.SetDefault(PostgresSSLMode, "disable")
}

func SetTestPostgresConfig() {
	viper.SetDefault(PostgresHost, "localhost")
	viper.SetDefault(PostgresPort, 5432)
	viper.SetDefault(PostgresDB, "judi_test_db")
	viper.SetDefault(PostgresUser, "judi_test")
	viper.SetDefault(PostgresPassword, "judi_test_pswd")
	viper.SetDefault(PostgresSSLMode, "disable")
}

// Redis

func SetDefaultRedisConfig() {
	viper.SetDefault(RedisHost, "localhost")
	viper.SetDefault(RedisPort, "6379")
	viper.SetDefault(RedisPassword, "judi_pswd")
}

func SetTestRedisConfig() {
	viper.SetDefault(RedisHost, "localhost")
	viper.SetDefault(RedisPort, "6379")
	viper.SetDefault(RedisPassword, "judi_test_pswd")
}

// Validation

func SetDefaultValidationConfig() {
	viper.SetDefault(MinUsernameLen, constants.MinUsernameLen)
	viper.SetDefault(MaxUsernameLen, constants.MaxUsernameLen)
}
