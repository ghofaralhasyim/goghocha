package configs

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/jet/v2"
)

func FiberConfig() fiber.Config {
	engine := jet.New("./templates", ".jet")

	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(60),
		Views:       engine,
	}
}
