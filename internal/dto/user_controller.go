package dto

import (
	"go-backend-template/internal/model"
	"go-backend-template/internal/util"
)

type UserController struct {
	userRepo   model.UserRepo
	userMapper model.OneWayUserMapper[User]
}

func NewUserController(userRepo model.UserRepo, userMapper model.OneWayUserMapper[User]) UserController {
	return UserController{
		userRepo:   userRepo,
		userMapper: userMapper,
	}
}

func (uc *UserController) GetUsers() ([]User, error) {
	users, err := uc.userRepo.GetUsers()

	if err != nil {
		return []User{}, err
	}

	return util.Map(users, func(user model.User) User {
		return uc.userMapper.FromUser(user)
	}), nil
}

func (uc *UserController) GetUsersByRole(role string) ([]User, error) {
	users, err := uc.userRepo.GetUsersByRole(role)

	if err != nil {
		return []User{}, err
	}

	return util.Map(users, func(user model.User) User {
		return uc.userMapper.FromUser(user)
	}), nil
}
