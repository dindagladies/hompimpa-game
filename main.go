package main

import (
	"hompimpa-game/config"
	"hompimpa-game/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()

	// app config
	app := fiber.New()
	route.SetupRoutes(app)

	// handle undefined routes
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Listen(":4000")
}
