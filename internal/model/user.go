package model

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Role         string
	FirstName    string
	LastName     string
}
