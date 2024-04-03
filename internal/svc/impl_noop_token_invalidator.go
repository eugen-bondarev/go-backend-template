package svc

import "time"

type NoopTokenInvalidator struct {
}

func NewNoopTokenInvalidator() ITokenInvalidatorSvc {
	return &NoopTokenInvalidator{}
}

func (ti *NoopTokenInvalidator) Invalidate(token string, until time.Time) {
}

func (ti *NoopTokenInvalidator) IsValid(token string) bool {
	return true
}
