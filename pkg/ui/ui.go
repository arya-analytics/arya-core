package ui

import "github.com/gofiber/fiber/v2"

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

const staticRoot = "./dist"

func (s *Server) BindTo(router fiber.Router) {
	router.Static("/", staticRoot)
}
