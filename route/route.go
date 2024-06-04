package route

import (
	"hompimpa-game/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/game", handler.GetGames)
	api.Get("/game/:code", handler.GetGameByCode)
	api.Post("/code", handler.CreateCode)
	api.Post("/player", handler.CreatePlayer)
	api.Post("/start", handler.Start)
	api.Post("/vote/:code/:playerId", handler.Vote)
	api.Get("/count/:code", handler.CountResult)

}
