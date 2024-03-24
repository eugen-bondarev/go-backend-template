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
		return err
	}

	encryptedPass := string(bytes)

	return authSvc.userRepo.CreateUser(email, encryptedPass, role)
}

func (authSvc *DefaultAuthSvc) AuthenticateUser(email, plainTextPassword string) (model.User, error) {
	user, err := authSvc.userRepo.GetUserByEmail(email)

	if err != nil {
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainTextPassword+authSvc.pepper))

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
