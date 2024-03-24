package model

type UserRepo interface {
	GetUserByEmail(email string) (User, error)
}
