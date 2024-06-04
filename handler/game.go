package handler

import (
	"hompimpa-game/config"
	"hompimpa-game/model"

	"github.com/gofiber/fiber/v2"
)

func GetGames(c *fiber.Ctx) error {
	var games []model.Game
	var db = config.DB
	db.Find(&games)

	if len(games) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No games found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Games found",
		"data":    games,
	})
}

func GetGameByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	round := c.Query("round")
	var games []model.Game
	var db = config.DB

	if round != "" {
		db.Where("code = ? AND round = ?", code, round).Order("round desc").Find(&games)
	} else {
		db.Where("code = ?", code).Order("round desc").Find(&games)
	}

	if len(games) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Game not found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Game found",
		"data":    games,
	})
}
