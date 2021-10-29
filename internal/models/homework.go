package models

import "time"

type Homework struct {
	ID        int64
	Name      string
	LectureID int64
	Filepath  string
	IsActive  bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

type HomeWorkWithFile struct {
	Homework
	File []byte
}
