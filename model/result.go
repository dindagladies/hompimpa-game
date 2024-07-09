package model

import (
	"time"
)

type Result struct {
	ID              int       `json:"id"`
	Code            string    `json:"code"`
	HandChoice      string    `json:"hand_choice"`
	WinnerPlayerIds []int     `json:"winner_player_ids" gorm:"serializer:json"`
	WinnerPlayer    []Player  `json:"winner_player" gorm:"many2many:player"`
	LoserPlayerIds  []int     `json:"loser_player_ids" gorm:"serializer:json"`
	LoserPlayer     []Player  `json:"loser_player" gorm:"many2many:player"`
	Round           int       `json:"round"`
	CreatedAt       time.Time `json:"created_at"`
}
