package auth

import (
	"github.com/arya-analytics/aryacore/pkg/api"
	"github.com/arya-analytics/aryacore/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	svc *auth.Service
}

const (
	cookieName    = "token"
	tokenPrefix   = "Bearer "
	groupEndpoint = "/auth"
	loginEndpoint = "/login"
)

func NewServer(svc *auth.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) BindTo(router fiber.Router) {
	r := router.Group(groupEndpoint)
	r.Post(loginEndpoint, s.login)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (s *Server) login(c *fiber.Ctx) error {
	var p loginRequest
	if err := c.BodyParser(&p); err != nil {
		return err
	}
	user, err := s.svc.Login(c.UserContext(), p.Username, p.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.ErrorResponse{
			Type:    api.ErrorTypeAuthentication,
			Message: "Invalid credentials.",
		})
	}
	token, err := auth.NewToken(user.ID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(api.ErrorResponse{Type: api.ErrorTypeUnknown, Message: "Internal server error."})
	}
	c.Cookie(&fiber.Cookie{Name: cookieName, Value: token})
	return c.JSON(loginResponse{Token: token})
}
