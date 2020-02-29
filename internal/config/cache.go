package config

// Cache holds all the configuration settings for the Redis client.
type Cache struct {
	Address                string `env:"REDIS_ADDRESS" default:"localhost"`
	Port                   string `env:"REDIS_PORT" default:":6379"`
	Password               string `env:"REDIS_PASSWORD" default:"Awesome85**"`
	DB                     int    `env:"REDIS_DB" default:"0"`
	PoolSize               int    `env:"REDIS_POOL_SIZE" default:"25"`
	MinimumIdleConnections int    `env:"REDIS_MINIMUM_IDLE_CONNECTIONS" default:"10"`
	MaximumRetries         int    `env:"REDIS_MAXIMUM_RETRIES" default:"1"`
	OnlineUserKeyFormat    string `env:"REDIS_ONLINE_USER_KEY_FORMAT" default:"online.user.%s"`
}
