package svc

import "time"

type ITokenInvalidator interface {
	Invalidate(token string, until time.Time)
	IsValid(token string) bool
}
