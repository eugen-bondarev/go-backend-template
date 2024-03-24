package impl

import "go-backend-template/internal/model"

type PGUserMapper struct {
}

func NewPGUserMapper() model.UserMapper[PGUser] {
	return &PGUserMapper{}
}

func (userMapper *PGUserMapper) FromUser(user model.User) PGUser {
	return PGUser{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}
}

func (userMapper *PGUserMapper) ToUser(pgUser PGUser) model.User {
	return model.User{
		ID:           pgUser.ID,
		Email:        pgUser.Email,
		PasswordHash: pgUser.PasswordHash,
		Role:         pgUser.Role,
	}
}
