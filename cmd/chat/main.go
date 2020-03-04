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
package main

import (
	"github.com/superhero-chat/cmd/chat/socketio"
	"github.com/superhero-chat/internal/config"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	socketIO, err := socketio.NewSocketIO(cfg)
	if err != nil {
		panic(err)
	}

	server, err := socketIO.NewSocketIOServer()
	if err != nil {
		panic(err)
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/", server)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServeTLS(
		cfg.App.Port,
		cfg.App.CertFile,
		cfg.App.KeyFile,
		serveMux,
	))
}
