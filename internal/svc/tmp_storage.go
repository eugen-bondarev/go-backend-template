package svc

import "time"

type ITmpStorage interface {
	Set(key string, value string, expiresAt time.Time) error
	Get(key string) (string, error)
}
