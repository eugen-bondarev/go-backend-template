package repo

import "go-backend-template/internal/model"

type IUserRepo interface {
	GetUserByEmail(email string) (model.User, error)
	GetUserByID(ID int) (model.User, error)
	GetUsers() ([]model.User, error)
	CreateUser(email, passwordHash, role string) error
	SetPasswordHashByEmail(email, passwordHash string) error
	DeleteUserByID(id int) error
}
