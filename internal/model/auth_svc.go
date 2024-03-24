package model

type AuthSvc interface {
	CreateUser(email, plainTextPassword, role string) error
	AuthenticateUser(email, plainTextPassword string) (User, error)
}
