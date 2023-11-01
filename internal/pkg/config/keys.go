package config

// Server
const (
	ServerName = "NAME"
	ServerPort = "PORT"
)

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

// S3
const (
	S3BucketName    = "S3_BUCKET_NAME"
	S3DefaultRegion = "S3_DEFAULT_REGION"
	S3Endpoint      = "S3_ENDPOINT"
)

// Validation
const (
	MinUsernameLen = "MIN_USERNAME_LEN"
	MaxUsernameLen = "MAX_USERNAME_LEN"
)
