package config

// Postgres
const (
	PostgresHost     = "PG_HOST"
	PostgresPort     = "PG_PORT"
	PostgresDB       = "PG_DB"
	PostgresUser     = "PG_USER"
	PostgresPassword = "PG_PASSWORD"
	PostgresSSLMode  = "PG_SSL_MODE"
)

// Redis
const (
	RedisHost     = "REDIS_HOST"
	RedisPort     = "REDIS_PORT"
	RedisPassword = "REDIS_PASSWORD"
)

// Validation
const (
	MinUsernameLen = "MIN_USERNAME_LEN"
	MaxUsernameLen = "MAX_USERNAME_LEN"
)
