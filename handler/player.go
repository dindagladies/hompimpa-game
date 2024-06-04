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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player created successfully",
		"data":    data,
	})
}

// TODO : get session player
