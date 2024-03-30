package impl

import (
	"go-backend-template/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type DefaultAuthSvc struct {
	userRepo model.UserRepo
	pepper   string
}

func NewDefaultAuthSvc(userRepo model.UserRepo, pepper string) model.AuthSvc {
	return &DefaultAuthSvc{
		userRepo: userRepo,
		pepper:   pepper,
	}
}

func (authSvc *DefaultAuthSvc) CreateUser(email, plainTextPassword, role string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword+authSvc.pepper), bcrypt.DefaultCost)

	if err != nil {
		return model.ErrAuthSvcCreateUserFailed
	}

	encryptedPass := string(bytes)

	err = authSvc.userRepo.CreateUser(email, encryptedPass, role)

	if err != nil {
		return model.ErrAuthSvcCreateUserFailed
	}

	return nil
}

func (authSvc *DefaultAuthSvc) SetPasswordByEmail(email, plainTextPassword string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword+authSvc.pepper), bcrypt.DefaultCost)

	if err != nil {
		return model.ErrAuthSvcCreateUserFailed
	}

	encryptedPass := string(bytes)

	err = authSvc.userRepo.SetPasswordHashByEmail(email, encryptedPass)

	return err
}

func (authSvc *DefaultAuthSvc) AuthenticateUser(email, plainTextPassword string) (model.User, error) {
	user, err := authSvc.userRepo.GetUserByEmail(email)

	if err != nil {
		return model.User{}, model.ErrAuthSvcAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainTextPassword+authSvc.pepper))

	if err != nil {
		return model.User{}, model.ErrAuthSvcAuthFailed
	}

	return user, nil
}
