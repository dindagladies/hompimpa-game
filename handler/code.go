package handler

import (
	"hompimpa-game/config"
	"hompimpa-game/model"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateCode(c *fiber.Ctx) error {
	data := new(model.Code)
	newCode := RandStringRunes(6)

	db := config.DB
	data.Code = newCode
	if err := db.Create(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create code",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Code created successfully",
		"data":    data,
	})
}
