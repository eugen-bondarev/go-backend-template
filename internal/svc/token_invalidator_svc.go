package svc

import "time"

type ITokenInvalidatorSvc interface {
	Invalidate(token string, until time.Time)
	IsValid(token string) bool
}
