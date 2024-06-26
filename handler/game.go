package handler

import (
	"hompimpa-game/config"
	"hompimpa-game/model"

	"github.com/gofiber/fiber/v2"
)

func GetGames(c *fiber.Ctx) error {
	var games []model.Game
	var db = config.DB
	db.Joins("Player").Find(&games)

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

	/*
	* Check user session
	 */
	userData, err := config.GetUserSession(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}

	if userData["ID"] == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "You are not logged in",
			"data":    nil,
		})
	}

	if db.Where("player_id = ?", userData["ID"].(int)).Find(&games).RowsAffected == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "You are not joining this game yet",
			"data":    nil,
		})
	}
	/* End */

	if round != "" && round != "next" {
		db.Joins("Player").Where("code = ? AND round = ?", code, round).Order("round desc").Find(&games)
	} else if round == "next" {
		db.Joins("Player").Where("code = ? AND round = ?", code, 0).Order("round desc").Find(&games)
	} else {
		db.Joins("Player").Where("code = ?", code).Order("round desc").Find(&games)
	}

	if len(games) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Game or round doesn't exist",
			"data":    nil,
		})
	}

	nextGameType := 1
	if len(games) <= 2 && round == "next" {
		nextGameType = 0
	} else if len(games) <= 2 && round == "" {
		nextGameType = 2
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Game found",
		"data":           games,
		"next_game_type": nextGameType,
		// "next_round":
	})
}
