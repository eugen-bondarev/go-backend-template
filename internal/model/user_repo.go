package model

type UserRepo interface {
	GetUserByEmail(email string) (User, error)
	CreateUser(email, passwordHash, role string) error
}