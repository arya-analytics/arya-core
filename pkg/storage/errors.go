package storage

import "fmt"

// |||| ERROR TYPES ||||

const (
	errKey = "storage"
)

type Error struct {
	Type ErrorType
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", errKey, e.Type)
}

func NewError(t ErrorType) Error {
	return Error{t}
}

type ErrorType int

//go:generate stringer -type=ErrorType
const (
	ErrTypeUniqueViolation ErrorType = iota
	ErrTypeRelationshipViolation
	ErrTypeNonPointer
	ErrTypeNonStructOrSlice
	ErrTypeInvalidField
	ErrTypeIncompatibleModels
	ErrTypeNoPK
)
