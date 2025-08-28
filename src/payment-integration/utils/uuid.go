package utils

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4 string
func GenerateUUID() string {
	return uuid.New().String()
}

// ValidateUUID checks if a string is a valid UUID
func ValidateUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
