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

	api.Get("/player", handler.GetPlayerLogin)
	api.Post("/player", handler.CreatePlayer)
	api.Post("/logout", handler.LogoutPlayer)
	api.Post("/exit/:code", handler.ExitGame)

	api.Post("/start", handler.Start)
	api.Get("/game/info/:code", handler.GetGameInfo)
	api.Post("/game/info/:code", handler.UpdateGameInfo)
	api.Post("/vote/:code/:playerId", handler.Vote)
	api.Get("/count/:code", handler.CountResult)
	api.Get("/result/:code", handler.GameResult)
}
