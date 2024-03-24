package model

type Payload struct {
	ID   int
	Role string
}

type SigningSvc interface {
	Sign(Payload) (string, error)
	Parse(string) (Payload, error)
}
