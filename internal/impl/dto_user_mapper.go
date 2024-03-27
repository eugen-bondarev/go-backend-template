package impl

import (
	"go-backend-template/internal/dto"
	"go-backend-template/internal/model"
)

type DTOUserMapper struct {
}

func NewDTOUserMapper() model.OneWayUserMapper[dto.User] {
	return &DTOUserMapper{}
}

func (userMapper *DTOUserMapper) FromUser(user model.User) dto.User {
	return dto.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
}
