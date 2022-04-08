package auth

import (
	"github.com/arya-analytics/aryacore/pkg/api"
	"github.com/arya-analytics/aryacore/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// |||| TOKEN ||||

func TokenMiddleware(c *fiber.Ctx) error {
	token, err := parseToken(c)
	if err != nil {
		return c.JSON(err)
	}
	if err = auth.ValidateToken(token); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(api.ErrorResponse{
			Type:    api.ErrorTypeUnauthorized,
			Message: "Invalid token",
		})
	}
	return c.Next()
}

func parseToken(c *fiber.Ctx) (string, error) {
	tokenParsers := []tokenParser{
		parseCookieToken,
		parseHeaderToken,
	}
	for _, tp := range tokenParsers {
		if token, err, ok := tp(c); ok {
			return token, err
		}
	}
	return "", api.ErrorResponse{
		Type:    api.ErrorTypeUnauthorized,
		Message: "No authentication token provided. Please provide token as cookie or in headers.",
	}
}

type tokenParser func(c *fiber.Ctx) (token string, err error, found bool)

func parseCookieToken(c *fiber.Ctx) (string, error, bool) {
	token := c.Cookies(cookieName)
	if len(token) == 0 {
		return "", nil, false
	}
	return token, nil, true
}

func parseHeaderToken(c *fiber.Ctx) (string, error, bool) {
	authHeader := c.Get("Authorization")
	if len(authHeader) == 0 {
		return "", nil, false
	}
	splitToken := strings.Split(authHeader, tokenPrefix)
	if len(splitToken) != 2 {
		return "", api.ErrorResponse{
			Type:    api.ErrorTypeUnauthorized,
			Message: "Invalid authorization header. Expected format: 'Authorization: Bearer <token>'.",
		}, true
	}
	return splitToken[1], nil, true
}
