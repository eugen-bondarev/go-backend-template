package model

type SigningSvc interface {
	Sign(ID int, role string) (string, error)
	Parse(string) (int, string, error)
}
