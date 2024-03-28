package util

import "errors"

type RequestError struct {
	StatusCode int
	Err        error
}

func NewRequestError(statusCode int, err error) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		Err:        err,
	}
}

func NewRequestErrorStr(statusCode int, str string) *RequestError {
	return NewRequestError(statusCode, errors.New(str))
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}
