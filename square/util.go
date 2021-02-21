package square

import (
	"fmt"

	"github.com/google/uuid"
)

// Generates a new temporary client ID for creating a new Square object.
func newTempID() *string {
	return strPtr(fmt.Sprint("#", uuid.New().String()))
}

// Returns a pointer to the specified string value.
func strPtr(value string) *string {
	return &value
}
