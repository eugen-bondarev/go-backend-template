package impl

import "go-backend-template/internal/model"

type PGUserRepo struct {
	pg         *Postgres
	userMapper model.UserMapper[PGUser]
}

func (userRepo *PGUserRepo) GetUserByEmail(email string) (model.User, error) {
	var users []PGUser

	err := userRepo.pg.GetDB().Select(&users, "SELECT * FROM users WHERE email = $1", email)

	if err != nil {
		return model.User{}, err
	}

	return userRepo.userMapper.ToUser(users[0]), nil
}

func NewPGUserRepo(pg *Postgres) model.UserRepo {
	return &PGUserRepo{
		pg:         pg,
		userMapper: NewPGUserMapper(),
	}
}
