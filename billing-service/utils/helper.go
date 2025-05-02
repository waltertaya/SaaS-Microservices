package utils

import "github.com/google/uuid"

func GenerateRef() string {
	id := uuid.New()

	return id.String()
}
