package models

import "time"

type Lecture struct {
	ID        int64
	Name      string
	Filepath  string
	CreatedAt time.Time
}

type LectureWithFile struct {
	Lecture
	File []byte
}
