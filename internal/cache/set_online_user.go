package cache

import (
	"fmt"
)

// SetOnlineUser stores user id of user that is currently has the app opened (is online) into Redis cache.
func (c *Cache) SetOnlineUser(userID string) error {
	err := c.Redis.Set(fmt.Sprintf(c.OnlineUserKeyFormat, userID), userID, 0, ).Err()
	if err != nil {
		return err
	}

	return nil
}