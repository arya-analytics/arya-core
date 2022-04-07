package errors

type ErrorResponse struct {
	Message string `json:"message"`
}

type ErrorType int

const (
	Unauthorized = 0
)