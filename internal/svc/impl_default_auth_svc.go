package svc

import (
	"go-backend-template/internal/model"
	"go-backend-template/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type DefaultAuthSvc struct {
	userRepo repo.IUserRepo
	pepper   string
}

func NewDefaultAuthSvc(userRepo repo.IUserRepo, pepper string) IAuthSvc {
	return &DefaultAuthSvc{
		userRepo: userRepo,
		pepper:   pepper,
	}
}

func (authSvc *DefaultAuthSvc) CreateUser(email, plainTextPassword, role string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword+authSvc.pepper), bcrypt.DefaultCost)

	if err != nil {
		return ErrAuthSvcCreateUserFailed
	}

	encryptedPass := string(bytes)

	err = authSvc.userRepo.CreateUser(email, encryptedPass, role)

	if err != nil {
		return ErrAuthSvcCreateUserFailed
	}

	return nil
}

func (authSvc *DefaultAuthSvc) SetPasswordByEmail(email, plainTextPassword string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword+authSvc.pepper), bcrypt.DefaultCost)

	if err != nil {
		return ErrAuthSvcCreateUserFailed
	}

	encryptedPass := string(bytes)

	err = authSvc.userRepo.SetPasswordHashByEmail(email, encryptedPass)

	return err
}

func (authSvc *DefaultAuthSvc) AuthenticateUser(email, plainTextPassword string) (model.User, error) {
	user, err := authSvc.userRepo.GetUserByEmail(email)

	if err != nil {
		return model.User{}, ErrAuthSvcAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainTextPassword+authSvc.pepper))

	if err != nil {
		return model.User{}, ErrAuthSvcAuthFailed
	}

	return user, nil
}
