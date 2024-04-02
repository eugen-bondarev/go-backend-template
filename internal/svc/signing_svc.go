package svc

import "time"

type ISigningSvc interface {
	Sign(claims map[string]any, expiration time.Time) (string, error)
	Parse(token string) (map[string]any, error)
}
