package repo

import (
	"errors"
	"go-backend-template/internal/model"
	"slices"
)

type MemUserRepo struct {
	users []model.User
}

func (userRepo *MemUserRepo) GetUserByEmail(email string) (model.User, error) {
	for _, user := range userRepo.users {
		if user.Email == email {
			return user, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func (userRepo *MemUserRepo) GetUserByID(ID int) (model.User, error) {
	for _, user := range userRepo.users {
		if user.ID == ID {
			return user, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func (userRepo *MemUserRepo) GetUsers() ([]model.User, error) {
	return userRepo.users, nil
}

func (userRepo *MemUserRepo) CreateUser(email, passwordHash, role string) error {
	userRepo.users = append(userRepo.users, model.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	})
	return nil
}

func (userRepo *MemUserRepo) SetPasswordHashByEmail(email, passwordHash string) error {
	for _, user := range userRepo.users {
		if user.Email == email {
			user.PasswordHash = passwordHash
			return nil
		}
	}
	return errors.New("user not found")
}

func (userRepo *MemUserRepo) DeleteUserByID(ID int) error {
	for i, user := range userRepo.users {
		if user.ID == ID {
			userRepo.users = slices.Delete(userRepo.users, i, i)
			return nil
		}
	}
	return nil
}

func NewMemUserRepo() IUserRepo {
	return &MemUserRepo{}
}
