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
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}
}

func (userMapper *PGUserMapper) ToUser(pgUser PGUser) model.User {
	return model.User{
		ID:           pgUser.ID,
		Username:     pgUser.Username,
		PasswordHash: pgUser.PasswordHash,
		Role:         pgUser.Role,
	}
}
