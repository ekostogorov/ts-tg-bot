package uid_helpers

import "github.com/google/uuid"

func GenerateUUIDv4() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return uid.String(), nil
}
