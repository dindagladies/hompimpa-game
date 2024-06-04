package handler

import (
	"hompimpa-game/config"
	"hompimpa-game/model"

	"github.com/gofiber/fiber/v2"
)

func Start(c *fiber.Ctx) error {
	game := new(model.Game)

	if err := c.BodyParser(game); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	var db = config.DB

	isPlayerAlreadyJoined := db.Where("player_id = ? AND code = ?", game.PlayerId, game.Code).First(&game).RowsAffected > 0

	if isPlayerAlreadyJoined {
		return c.Status(400).JSON(fiber.Map{
			"message": "Player already joined the game",
		})
	}

	if err := db.Create(&game).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Player cannot join the game",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player joined the game successfully",
		"data":    game,
	})
}

func Vote(c *fiber.Ctx) error {
	code := c.Params("code")
	playerId := c.Params("playerId")

	db := config.DB
	var game model.Game
	if err := db.Where("code = ? AND player_id = ?", code, playerId).First(&game).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Game not found",
		})
	}

	var UpdateVote model.UpdateVote
	if err := c.BodyParser(&UpdateVote); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	if err := db.Model(&game).Update("hand_choice", UpdateVote.HandChoice).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update vote",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Vote updated successfully",
		"data":    game,
	})
}

// TODO
func CountResult(c *fiber.Ctx) error {
	code := c.Params("code")
	var game model.Game
	var db = config.DB

	// get last round
	db.Where("code = ?", code).Order("round desc").First(&game)
	round := game.Round + 1

	// === if hompimpa game ====
	// get total vote
	var games []model.Game
	db.Where("code = ? AND round = ?", code, "NULL").Find(&games)

	// count vote
	var votes = make(map[string]int)
	for _, game := range games {
		_, ok := votes[game.HandChoice]
		if ok {
			votes[game.HandChoice]++
		} else {
			votes[game.HandChoice] = 1
		}
	}

	// get winner
	winner := 0
	handWinner := ""
	for key, value := range votes {
		if value > winner {
			winner = value
			handWinner = key
		}
	}

	// update result, round, and winner = last round
	db.Model(&game).Where("code = ? AND round = ? AND hand_choice NOT IN (?)", code, "NULL", handWinner).Update("round", round)
	db.Model(&game).Where("code = ? AND round = ? AND hand_choice = ?", code, "NULL", handWinner).Update("round", round+1)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Game result",
		"data":    votes,
		"winner":  handWinner,
	})
}
