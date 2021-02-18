package client

import (
	"fmt"

	"github.com/google/uuid"
)

// Generates a new idempotency key used when making Square API requests.
func newIdempotencyKey() *string {
	key := uuid.New().String()
	return &key
}

// Generates a new temporary ID used when creating new Square objects.
func newTempID() string {
	return fmt.Sprint("#", uuid.New().String())
}

func strPtr(value string) *string {
	return &value
}
