package util

import "errors"

type APIError struct {
	StatusCode int
	Err        error
}

func NewAPIError(statusCode int, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Err:        err,
	}
}

func NewAPIErrorStr(statusCode int, str string) *APIError {
	return NewAPIError(statusCode, errors.New(str))
}

func (r *APIError) Error() string {
	return r.Err.Error()
}
