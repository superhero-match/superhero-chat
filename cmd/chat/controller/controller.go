package controller

import (
	"net/http"

	"github.com/superhero-chat/cmd/chat/service"
	"github.com/superhero-chat/internal/config"
)

// Controller holds the Controller data.
type Controller struct {
	Service *service.Service
}

// NewController returns new controller.
func NewController(cfg *config.Config) (*Controller, error) {
	srv, err := service.NewService(cfg)
	if err != nil {
		return nil, err
	}

	return &Controller{
		Service: srv,
	}, nil
}

// SetupRoutes configures endpoints.
func (c *Controller) SetupRoutes() {
	http.HandleFunc("/ws", WsEndpoint)
}