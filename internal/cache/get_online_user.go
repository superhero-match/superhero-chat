package cache

import (
	"github.com/go-redis/redis"
)

// GetOnlineUser fetches online user id from cache.
func (c *Cache) GetOnlineUser(key string) (string, error) {
	onlineUserID, err := c.Redis.Get(key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	if len(onlineUserID) == 0 {
		return "", nil
	}

	return onlineUserID, nil
}
