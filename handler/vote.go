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
	if err := db.Where("code = ? AND player_id = ? AND round = ?", code, playerId, 0).First(&game).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Game or player doesn't exist in this round",
		})
	}

	var UpdateVote model.UpdateVote
	if err := c.BodyParser(&UpdateVote); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	if err := db.Model(&game).Updates(model.Game{HandChoice: UpdateVote.HandChoice, GameTypeId: UpdateVote.GameTypeId}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update vote",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Vote updated successfully",
		"data":    game,
	})
}

func CountResult(c *fiber.Ctx) error {
	code := c.Params("code")
	var game model.Game
	var db = config.DB

	// get last round
	db.Where("code = ?", code).Order("round desc").First(&game)
	newRound := game.Round + 1

	// get votes
	var games []model.Game
	db.Where("code = ? AND round = ?", code, 0).Find(&games)

	if len(games) == 0 && game.Round > 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Game has finished",
		})
	} else if len(games) == 0 && game.Round == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Game not started yet",
		})
	}

	if len(games) == 2 { //game type id = 2
		handFirstPlayer := games[0].HandChoice
		handSecondPlayer := games[1].HandChoice
		winner := ""

		if (handFirstPlayer == "rock" && handSecondPlayer == "scissors") || (handFirstPlayer == "scissors" && handSecondPlayer == "rock") {
			winner = "rock"
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "scissors").Update("round", newRound)
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "rock").Update("round", newRound+1)
		} else if handFirstPlayer == "rock" && handSecondPlayer == "paper" || handFirstPlayer == "paper" && handSecondPlayer == "rock" {
			winner = "paper"
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "rock").Update("round", newRound)
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "paper").Update("round", newRound+1)
		} else if (handFirstPlayer == "scissors" && handSecondPlayer == "paper") || handFirstPlayer == "paper" && handSecondPlayer == "scissors" {
			winner = "scissors"
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "paper").Update("round", newRound)
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "scissors").Update("round", newRound+1)
		} else if handFirstPlayer == handSecondPlayer {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"continue_round": true,
				"message":        "Game result",
				"next_game_type": 2,
				"vote_result": map[string]int{
					handFirstPlayer:  1,
					handSecondPlayer: 1,
				},
				"winner": "draw",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"continue_round": false,
			"message":        "Game result",
			"next_game_type": 0,
			"vote_result": map[string]int{
				handFirstPlayer:  1,
				handSecondPlayer: 1,
			},
			"winner": winner,
		})
	}

	var votes = make(map[string]int)
	for _, game := range games {
		_, exist := votes[game.HandChoice]
		if exist {
			votes[game.HandChoice]++
		} else {
			votes[game.HandChoice] = 1
		}
	}

	if len(votes) <= 1 {
		return c.Status(400).JSON(fiber.Map{
			"message":     "Game not finished yet",
			"vote_result": votes,
		})
	}

	// get winner from votes
	winnerCount := 0
	winnerChoice := ""
	for hand, vote := range votes {
		if winnerCount < vote {
			winnerCount = vote
			winnerChoice = hand
		} else if winnerCount == vote {
			winnerChoice = "draw"
		}
	}

	// update the looser player
	if winnerChoice != "draw" && winnerChoice != "" {
		db.Model(&games).Where("code = ? AND round = ? AND hand_choice NOT IN (?)", code, 0, winnerChoice).Update("round", newRound)
		// db.Model(&game).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, handWinner).Update("round", 0)
	}

	// if winner has more than 1
	game_type_id := 1
	if winnerCount > 1 {
		game_type_id = 2
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"continue_round": winnerCount != 1,
		"message":        "Game result",
		"next_game_type": game_type_id,
		"vote_result":    votes,
		"winner":         winnerChoice,
	})
}
