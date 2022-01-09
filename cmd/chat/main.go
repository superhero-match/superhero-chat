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
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/superhero-match/superhero-chat/cmd/chat/socketio"
	"github.com/superhero-match/superhero-chat/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	router := gin.New()

	socketIO, err := socketio.NewSocketIO(cfg)
	if err != nil {
		panic(err)
	}

	server, err := socketIO.NewSocketIOServer()
	if err != nil {
		panic(err)
	}

	go server.Serve()
	defer server.Close()

	router.GET("/*any", gin.WrapH(server))
	router.POST("/*any", gin.WrapH(server))

	err = router.Run(cfg.App.Port)
	if err != nil {
		panic(err)
	}
}

// Health is used for health checks from loadbalancer.
func Health(c *gin.Context) {
	c.Status(http.StatusOK)
}
