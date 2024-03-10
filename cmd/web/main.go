package main

import (
	"time"

	"github.com/ghofaralhasyim/goghocha/pkg/configs"
	"github.com/ghofaralhasyim/goghocha/pkg/handlers"
	"github.com/ghofaralhasyim/goghocha/pkg/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

func main() {
	config := configs.FiberConfig()

	app := fiber.New(config)

	routes.PublicRoutes(app)

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
	}))

	app.Static("/static", "./static")

	go handlers.ListenToWs()

	app.Listen(":3000")
}
