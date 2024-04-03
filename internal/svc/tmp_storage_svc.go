package svc

import "time"

type ITmpStorageSvc interface {
	Set(key string, value string, expiresAt time.Time) error
	Get(key string) (string, error)
}
