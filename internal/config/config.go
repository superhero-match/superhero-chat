package config

import (
	"github.com/jinzhu/configor"
)

// Config holds the configuration.
type Config struct {
	DB       *DB
	RabbitMQ *RabbitMQ
}

// NewConfig returns the configuration.
func NewConfig() (cnf *Config, e error) {
	var cfg Config

	if err := configor.Load(&cfg, "config.yml"); err != nil {
		return nil, err
	}

	return &cfg, nil
}
