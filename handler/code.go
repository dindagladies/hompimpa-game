package handler

import (
	"hompimpa-game/config"
	"hompimpa-game/model"
	"math/rand"
	"time"

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

	/* Get session data */
	userData, err := config.GetUserSession(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}

	var player model.Player
	db := config.DB
	db.Find(&player, userData["ID"].(int))

	if player.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Please login first",
			"data":    nil,
		})
	}
	/* End */

	data.Code = newCode
	data.HostID = player.ID

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

func UpdateGameInfo(c *fiber.Ctx) error {
	code := c.Params("code")
	db := config.DB
	var data model.Code

	if err := db.Where("code = ?", code).First(&data).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Game not found",
		})
	}

	if db.Where("code = ?", code).First((&data)).RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Game not found",
		})
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	data.StartedAt = time.Now().In(loc).Format("2006-01-02 15:04:05")

	if err := db.Where("code = ?", code).Updates(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update game",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Game updated successfully",
		"data":    data,
	})
}
