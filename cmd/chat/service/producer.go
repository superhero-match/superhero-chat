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
package service

import (
	"github.com/superhero-match/superhero-chat/cmd/chat/model"
	"github.com/superhero-match/superhero-chat/cmd/chat/service/mapper"
)

// StoreMessage publishes new message on Kafka topic for it to be
// consumed by consumer and stored in Cache and DB.
func (s *Service) StoreMessage(m model.Message, isOnline bool, createdAt string) error {
	return s.Producer.StoreMessage(mapper.MapAPIMessageToProducer(m, isOnline, createdAt))
}