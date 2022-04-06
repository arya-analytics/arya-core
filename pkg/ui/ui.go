package ui

import "github.com/gofiber/fiber/v2"

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

const distRoot = "dist"

func (s *Server) BindTo(router fiber.Router) {
	router.Static("/", distRoot)
}
