package svc

import "time"

type Token struct {
	Value     string
	ExpiresAt time.Time
}

type ISigning interface {
	Sign(claims map[string]any, expiration time.Time) (Token, error)
	Parse(token string) (map[string]any, error)
}
