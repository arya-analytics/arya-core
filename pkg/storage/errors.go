package storage

import "fmt"

// |||| ERROR TYPES ||||

const (
	errKey = "storage"
)

type Error struct {
	Base    error
	Type    ErrorType
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s - %s", errKey, e.Type, e.Message)
}

type ErrorType int

//go:generate stringer -type=ErrorType
const (
	ErrTypeUnknown ErrorType = iota
	ErrTypeItemNotFound
	ErrTypeUniqueViolation
	ErrTypeRelationshipViolation
	ErrTypeInvalidField
	ErrTypeNoPK
	ErrTypeMigration
	ErrTypeInvalidArgs
)
