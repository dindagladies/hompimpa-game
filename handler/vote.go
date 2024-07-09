package handler

import (
	"hompimpa-game/config"
	"hompimpa-game/model"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Start(c *fiber.Ctx) error {
	// TODO: implement session player
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

	db.Joins("Player").Where("code = ? AND player_id = ?", code, playerId).First(&game)

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
		var playerData = make(map[string]int)
		playerData[handFirstPlayer] = (games[0].PlayerId)
		playerData[handSecondPlayer] = (games[1].PlayerId)
		log.Println(playerData)
		winner := ""
		var winnerPlayerId int
		var looserPlayerId int

		isFirstPlayerDisqualification := handFirstPlayer == ""
		isSecondPlayerDisqualification := handSecondPlayer == ""

		/* Check isqualification */
		if (isFirstPlayerDisqualification) || (isSecondPlayerDisqualification) {
			var winnerPlayer []int
			var looserPlayer []int
			if (isFirstPlayerDisqualification) && (isSecondPlayerDisqualification) {
				winnerPlayer = []int{0}
				looserPlayer = []int{games[0].PlayerId, games[1].PlayerId}
				db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "").Update("round", newRound)
			} else if isFirstPlayerDisqualification {
				winnerPlayer = []int{games[1].PlayerId}
				looserPlayer = []int{games[0].PlayerId}
				db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "").Update("round", newRound)
				db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, games[1].HandChoice).Update("round", newRound+1)
			} else if isSecondPlayerDisqualification {
				winnerPlayer = []int{games[0].PlayerId}
				looserPlayer = []int{games[1].PlayerId}
				db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "").Update("round", newRound)
				db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, games[0].HandChoice).Update("round", newRound+1)
			}

			result := model.Result{
				Code:            code,
				HandChoice:      "DRAW",
				WinnerPlayerIds: winnerPlayer,
				LoserPlayerIds:  looserPlayer,
				Round:           newRound,
			}

			log.Println(result)

			if err := db.Create(&result).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Failed to insert result",
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
				"winner": "",
			})
		}
		/* End Check isqualification */

		if (handFirstPlayer == "rock" && handSecondPlayer == "scissors") || (handFirstPlayer == "scissors" && handSecondPlayer == "rock") {
			winner = "rock"
			winnerPlayerId = playerData["rock"]
			looserPlayerId = playerData["scissors"]
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "scissors").Update("round", newRound)
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "rock").Update("round", newRound+1)

		} else if handFirstPlayer == "rock" && handSecondPlayer == "paper" || handFirstPlayer == "paper" && handSecondPlayer == "rock" {
			winner = "paper"
			winnerPlayerId = playerData["paper"]
			looserPlayerId = playerData["rock"]
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "rock").Update("round", newRound)
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "paper").Update("round", newRound+1)
		} else if (handFirstPlayer == "scissors" && handSecondPlayer == "paper") || handFirstPlayer == "paper" && handSecondPlayer == "scissors" {
			winner = "scissors"
			winnerPlayerId = playerData["scissors"]
			looserPlayerId = playerData["paper"]
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "paper").Update("round", newRound)
			db.Model(&games).Where("code = ? AND round = ? AND hand_choice = ?", code, 0, "scissors").Update("round", newRound+1)
		} else if handFirstPlayer == handSecondPlayer {
			// TODO: test this case (draw in game type 2)
			result := model.Result{
				Code:            code,
				HandChoice:      "DRAW",
				WinnerPlayerIds: []int{games[0].PlayerId, games[1].PlayerId},
				LoserPlayerIds:  []int{},
				Round:           newRound,
			}

			log.Println(result)

			if err := db.Create(&result).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Failed to insert result",
				})
			}

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
		} else {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid hand choice",
			})
		}

		// TODO: update data in code table (is_finished)
		var codeData model.Code
		if err := db.Model(&codeData).Where("code = ?", code).Update("is_finished", true).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to update game finished",
			})
		}

		result := model.Result{
			Code:            code,
			HandChoice:      winner,
			WinnerPlayerIds: []int{winnerPlayerId},
			LoserPlayerIds:  []int{looserPlayerId},
			Round:           newRound,
		}

		if err := db.Create(&result).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to insert result",
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
	var listPlayerId []int
	for _, game := range games {
		_, exist := votes[game.HandChoice]
		if exist {
			votes[game.HandChoice]++
		} else {
			votes[game.HandChoice] = 1
		}
		listPlayerId = append(listPlayerId, game.PlayerId)
	}

	if len(votes) <= 1 {
		if len(listPlayerId) > 1 {
			handChoice := ""
			total := 0
			for hand, vote := range votes {
				handChoice = hand
				total = vote
			}

			result := model.Result{
				Code:            code,
				HandChoice:      "DRAW",
				WinnerPlayerIds: listPlayerId,
				LoserPlayerIds:  []int{},
				Round:           newRound,
			}

			log.Println(result)

			if err := db.Create(&result).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{
					"message": "Failed to insert result",
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"continue_round": true,
				"message":        "Game result",
				"next_game_type": 2,
				"vote_result": map[string]int{
					handChoice: total,
				},
				"winner": "draw",
			})
		}

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

		// insert result
		var winnerPlayerIds []int
		var loserPlayerIds []int
		for _, game := range games {
			if game.HandChoice == winnerChoice {
				winnerPlayerIds = append(winnerPlayerIds, game.PlayerId)
			} else {
				loserPlayerIds = append(loserPlayerIds, game.PlayerId)
			}
		}

		log.Println(winnerPlayerIds, loserPlayerIds)

		result := model.Result{
			Code:            code,
			HandChoice:      winnerChoice,
			WinnerPlayerIds: winnerPlayerIds,
			LoserPlayerIds:  loserPlayerIds,
			Round:           newRound,
		}

		log.Println(result)

		if err := db.Create(&result).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to insert result",
			})
		}
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
