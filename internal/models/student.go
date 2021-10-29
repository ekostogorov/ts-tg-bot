package models

import (
	"time"
)

type Student struct {
	ID             int64
	Name           string
	Login          string
	Folder         string
	TelegramUserID string
	IsActivated    bool
	ActivatedAt    time.Time
	CreatedAt      time.Time
}
