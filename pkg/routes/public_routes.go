package routes

import (
	"github.com/ghofaralhasyim/goghocha/pkg/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {

	route := a

	route.Get("/", handlers.Homepage)
	route.Get("/gocha", websocket.New(handlers.WsEndpoint))
}
