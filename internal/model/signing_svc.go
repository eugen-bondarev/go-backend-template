package model

type SigningSvc interface {
	Sign(map[string]any) (string, error)
	Parse(string) (map[string]any, error)
}
