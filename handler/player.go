package handler

import (
	"hompimpa-game/config"
	"hompimpa-game/model"

	"github.com/gofiber/fiber/v2"
)

func CreatePlayer(c *fiber.Ctx) error {
	data := new(model.Player)
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Failed to parse player data",
		})
	}

	db := config.DB
	if err := db.Create(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create player",
		})
	}

	if err := config.CreateUserSession(c, data.ID, data.Username); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create session player",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player created successfully",
		"data":    data,
	})
}

func GetPlayerLogin(c *fiber.Ctx) error {
	userData, err := config.GetUserSession(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}

	var players []model.Player
	var db = config.DB
	db.Find(&players, userData["ID"].(int))

	if len(players) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No players found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Players found",
		"id":       players[0].ID,
		"username": players[0].Username,
	})
}
