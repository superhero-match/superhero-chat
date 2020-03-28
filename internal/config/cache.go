/*
  Copyright (C) 2019 - 2020 MWSOFT
  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.
  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.
  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package config

// Cache holds all the configuration settings for the Redis client.
type Cache struct {
	Address                string `env:"REDIS_ADDRESS" default:"192.168.0.105"`
	Port                   string `env:"REDIS_PORT" default:":6379"`
	Password               string `env:"REDIS_PASSWORD" default:"Awesome85**"`
	DB                     int    `env:"REDIS_DB" default:"0"`
	PoolSize               int    `env:"REDIS_POOL_SIZE" default:"25"`
	MinimumIdleConnections int    `env:"REDIS_MINIMUM_IDLE_CONNECTIONS" default:"10"`
	MaximumRetries         int    `env:"REDIS_MAXIMUM_RETRIES" default:"1"`
	OnlineUserKeyFormat    string `env:"REDIS_ONLINE_USER_KEY_FORMAT" default:"online.user.%s"`
}
