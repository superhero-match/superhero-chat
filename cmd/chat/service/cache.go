package service

import (
	"github.com/go-redis/redis"
)

// SetOnlineUser stores user id of user that is currently has the app opened (is online) into Redis cache.
func (s *Service) SetOnlineUser(userID string) error {
	err := s.Cache.SetOnlineUser(userID)
	if err != nil {
		return err
	}

	return nil
}

// GetOnlineUser fetches online user id from cache.
func (s *Service) GetOnlineUser(key string) (string, error) {
	onlineUserID, err := s.Cache.GetOnlineUser(key)
	if err != nil && err != redis.Nil {
		return "", err
	}

	if len(onlineUserID) == 0 {
		return "", nil
	}

	return onlineUserID, nil
}
