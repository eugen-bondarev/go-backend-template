package model

type UserRepo interface {
	GetUserByEmail(email string) (User, error)
	GetUserByID(ID int) (User, error)
	GetUsers() ([]User, error)
	CreateUser(email, passwordHash, role string) error
	SetPasswordHashByEmail(email, passwordHash string) error
	DeleteUserByID(id int) error
}
