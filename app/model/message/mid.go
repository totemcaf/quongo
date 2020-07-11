package message

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

// MID is the message identifier
type MID string

// Empty is the empty MID, it is a marker, it is not e valid one
const Empty = MID("")

// Valid ID are formed only with letters, digits, and specific separators: ".-:_"
var validID = regexp.MustCompile(`^[-.:_A-Za-z0-9]+$`)

// ParseID checks the string is a valid message ID and convert to it
func ParseID(mid string) (MID, error) {
	if mid == "" {
		return Empty, errors.New("Invalid empty Messgae ID")
	}

	if !validID.MatchString(mid) {
		return Empty, fmt.Errorf("Invalid ID '%v'", mid)
	}

	return MID(mid), nil
}

// NewID creates a reandom message ID universally unique
func NewID() MID {
	return MID(uuid.New().String())
}

// IsEmpty checks if the provided mid is the empty one
func (id MID) IsEmpty() bool {
	return id == Empty
}

// ToString returns the message id as a string
func (id MID) ToString() string {
	return string(id)
}
