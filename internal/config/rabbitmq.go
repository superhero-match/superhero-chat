package config

// RabbitMQ holds the configuration values for the RabbitMQ.
type RabbitMQ struct {
	Host     string `env:"RABBITMQ_HOST" default:"amqp://%s:%s@%s:%d/"`
	User     string `env:"RABBITMQ_USER" default:"superhero"`
	Password string `env:"RABBITMQ_PASSWORD" default:"match"`
	Address  string `env:"RABBITMQ_ADDRESS" default:"localhost"`
	Port     int    `env:"RABBITMQ_PORT" default:"5672"`
}
