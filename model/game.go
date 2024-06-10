package model

import (
	"time"
)

type Game struct {
	ID         int       `json:"id"`
	Code       string    `json:"code"`
	Round      int       `json:"round"`
	PlayerId   int       `json:"player_id"`
	Player     Player    `json:"player" gorm:"foreignKey:player_id"`
	GameTypeId int       `json:"game_type_id"`
	HandChoice string    `json:"hand_choice"`
	CreatedAt  time.Time `json:"created_at"`
}

type UpdateVote struct {
	HandChoice string `json:"hand_choice"`
	GameTypeId int    `json:"game_type_id"`
}
