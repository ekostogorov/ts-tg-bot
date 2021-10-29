package types

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrUserExists        = errors.New("user is already registered")
	ErrGetLectrureNumber = errors.New("failed to get lecture number from message")
	ErrHWExpired         = errors.New("homework expired")
)
