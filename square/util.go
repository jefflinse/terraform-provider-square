package square

import (
	"fmt"

	"github.com/google/uuid"
)

// Generates a new idempotency key for a Square API request.
func newIdempotencyKey() *string {
	key := uuid.New().String()
	return &key
}

// Generates a new temporary client ID for creating a new Square object.
func newTempID() string {
	return fmt.Sprint("#", uuid.New().String())
}

// Returns a pointer to the specified string value.
func strPtr(value string) *string {
	return &value
}
