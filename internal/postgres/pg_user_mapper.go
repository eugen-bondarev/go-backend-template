package postgres

import "go-backend-template/internal/model"

type PGUserMapper struct {
}

func NewPGUserMapper() model.ModelMapper[model.User, PGUser] {
	return &PGUserMapper{}
}

func (userMapper *PGUserMapper) FromModel(user model.User) PGUser {
	return PGUser{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}
}

func (userMapper *PGUserMapper) ToModel(pgUser PGUser) model.User {
	return model.User{
		ID:           pgUser.ID,
		Email:        pgUser.Email,
		PasswordHash: pgUser.PasswordHash,
		Role:         pgUser.Role,
		FirstName:    pgUser.FirstName,
		LastName:     pgUser.LastName,
	}
}
