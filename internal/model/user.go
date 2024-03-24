package model

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
}
