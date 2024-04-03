package svc

import "time"

type DefaultTokenInvalidator struct {
	tmpStorageSvc ITmpStorageSvc
}

func NewDefaultTokenInvalidator(tmpStorageSvc ITmpStorageSvc) ITokenInvalidatorSvc {
	return &DefaultTokenInvalidator{
		tmpStorageSvc: tmpStorageSvc,
	}
}

func (ti *DefaultTokenInvalidator) Invalidate(token string, until time.Time) {
	ti.tmpStorageSvc.Set(token, "1", until)
}

func (ti *DefaultTokenInvalidator) IsValid(token string) bool {
	val, err := ti.tmpStorageSvc.Get(token)

	if err != nil {
		return true
	}

	if val == "" {
		return true
	}

	return false
}
