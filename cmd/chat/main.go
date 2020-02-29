package main

import (
	"github.com/superhero-chat/cmd/chat/controller"
	"github.com/superhero-chat/internal/config"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	ctrl, err := controller.NewController(cfg)
	if err != nil {
		panic(err)
	}

	ctrl.SetupRoutes()

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServeTLS(
		cfg.App.Port,
		cfg.App.CertFile,
		cfg.App.KeyFile,
		nil,
	))
}
