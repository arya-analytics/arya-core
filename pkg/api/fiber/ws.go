package fiber

import (
	"github.com/arya-analytics/aryacore/pkg/ws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebsocketHandler(f func(c *ws.Conn)) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) { f(&ws.Conn{Conn: c}) })
}
