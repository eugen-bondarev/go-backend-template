package model

import "time"

type SigningSvc interface {
	Sign(claims map[string]any, expiration time.Time) (string, error)
	Parse(token string) (map[string]any, error)
}
