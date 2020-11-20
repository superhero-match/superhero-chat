/*
  Copyright (C) 2019 - 2021 MWSOFT
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
	"github.com/gin-gonic/gin"
	"github.com/superhero-match/superhero-chat/cmd/chat/socketio"
	"github.com/superhero-match/superhero-chat/internal/config"
	"github.com/superhero-match/superhero-chat/internal/health"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	client := health.NewClient(cfg)

	router := gin.New()

	socketIO, err := socketio.NewSocketIO(cfg)
	if err != nil {
		_ = client.ShutdownHealthServer()

		panic(err)
	}

	server, err := socketIO.NewSocketIOServer()
	if err != nil {
		_ = client.ShutdownHealthServer()

		panic(err)
	}

	go server.Serve()
	defer server.Close()

	router.GET("/*any", gin.WrapH(server))
	router.POST("/*any", gin.WrapH(server))

	err = router.RunTLS(
		cfg.App.Port,
		cfg.App.CertFile,
		cfg.App.KeyFile,
	)
	if err != nil {
		_ = client.ShutdownHealthServer()

		panic(err)
	}

	_ = client.ShutdownHealthServer()
}
