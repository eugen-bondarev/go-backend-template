package model

type UserRepo interface {
	GetUserByID(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUsers() ([]User, error)
	CreateUser(email, passwordHash, role string) error
	DeleteUserByID(id int) error
}
