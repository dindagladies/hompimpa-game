package model

import "time"

type Code struct {
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	IsFinished bool      `json:"is_finished"`
}
