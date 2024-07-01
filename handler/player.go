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
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"players_username_key\" (SQLSTATE 23505)" {
			return c.Status(400).JSON(fiber.Map{
				"message": "Username already exists",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create player.",
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
	/* Get session data */
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
		return c.Status(404).JSON(fiber.Map{
			"message": "Please login first",
			"data":    nil,
		})
	}
	/* End */

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Players found",
		"id":       players[0].ID,
		"username": players[0].Username,
	})
}

func ExitGame(c *fiber.Ctx) error {
	userData, err := config.GetUserSession(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get session",
			"data":    nil,
		})
	}

	var player model.Player
	var db = config.DB
	db.Find(&player, userData["ID"].(int))
	if player.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Player not found",
		})
	}

	code := c.Params("code")
	var game model.Game
	db.Where("code = ? AND player_id = ?", code, player.ID).Order("id desc").First(&game)
	if game.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Player not found in this game",
		})
	}

	if err := db.Delete(&game).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to exit game",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player exit game successfully",
	})
}

func LogoutPlayer(c *fiber.Ctx) error {
	_, err := config.GetUserSession(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}

	sessionRemove, err := config.RemoveUserSession(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to logout",
		})
	}

	if sessionRemove {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logout success",
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": "Failed to logout",
	})
}
