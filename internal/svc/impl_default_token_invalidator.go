package svc

import "time"

type DefaultTokenInvalidator struct {
	tmpStorage ITmpStorage
}

func NewDefaultTokenInvalidator(tmpStorage ITmpStorage) ITokenInvalidator {
	return &DefaultTokenInvalidator{
		tmpStorage: tmpStorage,
	}
}

func (ti *DefaultTokenInvalidator) Invalidate(token string, until time.Time) {
	ti.tmpStorage.Set(token, "1", until)
}

func (ti *DefaultTokenInvalidator) IsValid(token string) bool {
	val, err := ti.tmpStorage.Get(token)

	if err != nil {
		return true
	}

	if val == "" {
		return true
	}

	return false
}
