package svc

import (
	"go-backend-template/internal/model"
	"go-backend-template/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type DefaultAuth struct {
	userRepo repo.IUserRepo
	pepper   string
}

func NewDefaultAuth(userRepo repo.IUserRepo, pepper string) IAuth {
	return &DefaultAuth{
		userRepo: userRepo,
		pepper:   pepper,
	}
}

func (auth *DefaultAuth) CreateUser(email, plainTextPassword, role string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword+auth.pepper), bcrypt.DefaultCost)

	if err != nil {
		return ErrAuthCreateUserFailed
	}

	encryptedPass := string(bytes)

	err = auth.userRepo.CreateUser(email, encryptedPass, role)

	if err != nil {
		return ErrAuthCreateUserFailed
	}

	return nil
}

func (auth *DefaultAuth) SetPasswordByEmail(email, plainTextPassword string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword+auth.pepper), bcrypt.DefaultCost)

	if err != nil {
		return ErrAuthCreateUserFailed
	}

	encryptedPass := string(bytes)

	err = auth.userRepo.SetPasswordHashByEmail(email, encryptedPass)

	return err
}

func (auth *DefaultAuth) AuthenticateUser(email, plainTextPassword string) (model.User, error) {
	user, err := auth.userRepo.GetUserByEmail(email)

	if err != nil {
		return model.User{}, ErrAuthAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainTextPassword+auth.pepper))

	if err != nil {
		return model.User{}, ErrAuthAuthFailed
	}

	return user, nil
}
