package api

import "fmt"

type ErrorResponse struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s - %s", e.Type, e.Message)
}

//go:generate stringer -type=ErrorType
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeUnauthorized
	ErrorTypeAuthentication
	ErrorTypeInvalidArguments
)

func NewErrorResponse(t ErrorType, message string) ErrorResponse {
	return ErrorResponse{
		Type:    t,
		Message: message,
	}
}
