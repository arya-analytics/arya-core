package ui

import (
	"github.com/gofiber/fiber/v2"
)

// Server servers the UI and assets related to it.
type Server struct{}

// NewServer creates a new server.
func NewServer() *Server {
	return &Server{}
}

func (s *Server) BindTo(router fiber.Router) {
	router.Use(distHandler())
}
