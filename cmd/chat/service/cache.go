/*
  Copyright (C) 2019 - 2022 MWSOFT
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
package service

import (
	"github.com/go-redis/redis"
)

// SetOnlineUser stores user id of user that is currently has the app opened (is online) into Redis cache.
func (s *service) SetOnlineUser(key string, userID string) error {
	err := s.Cache.SetOnlineUser(key, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetOnlineUser fetches online user id from cache.
func (s *service) GetOnlineUser(key string) (string, error) {
	onlineUserID, err := s.Cache.GetOnlineUser(key)
	if err != nil && err != redis.Nil {
		return "", err
	}

	if len(onlineUserID) == 0 {
		return "", nil
	}

	return onlineUserID, nil
}

// DeleteOnlineUser deletes online user form Redis cache when user disconnects.
func (s *service) DeleteOnlineUser(keys []string, userID string) error {
	return s.Cache.DeleteOnlineUser(keys, userID)
}
