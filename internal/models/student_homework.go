package models

import (
	"time"
)

type StudentHomework struct {
	StudentID  int64
	HomeworkID int64
	PassedAt   time.Time
	Filepath   string
	Grade      int64
}
