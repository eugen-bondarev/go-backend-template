package util

import (
	"go-backend-template/internal/localization"
)

type APIError struct {
	StatusCode int
	Message    localization.Message
}

func NewAPIError(statusCode int, message localization.Message) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (r *APIError) Error() string {
	return r.Message.GetContent()
}
