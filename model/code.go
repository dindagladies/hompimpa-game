package model

import "time"

type Code struct {
	Code              string    `json:"code"`
	CreatedAt         time.Time `json:"created_at"`
	IsFinished        bool      `json:"is_finished"`
	HostID            int       `json:"host_id"`
	TotalActivePlayer int       `json:"active_player" gorm:"-"`
	StartedAt         string    `json:"started_at"`
}
