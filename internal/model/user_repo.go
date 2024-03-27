package model

type UserRepo interface {
	GetUserByEmail(email string) (User, error)
	GetUsers() ([]User, error)
	GetUsersByRole(role string) ([]User, error)
	CreateUser(email, passwordHash, role string) error
	DeleteUserByID(id int) error
}
