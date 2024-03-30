package impl

import (
	"errors"
	"go-backend-template/internal/model"

	"github.com/eugen-bondarev/go-slice-helpers/parallel"
)

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

	if len(users) == 0 {
		return model.User{}, errors.New("user not found")
	}

	return userRepo.userMapper.ToUser(users[0]), nil
}

func (userRepo *PGUserRepo) GetUsers() ([]model.User, error) {
	var users []PGUser

	err := userRepo.pg.GetDB().Select(&users, "SELECT * FROM users")

	if err != nil {
		return []model.User{}, err
	}

	if len(users) == 0 {
		return []model.User{}, errors.New("user not found")
	}

	return parallel.Map(users, userRepo.userMapper.ToUser), nil
}

func (userRepo *PGUserRepo) CreateUser(email, passwordHash, role string) error {
	_, err := userRepo.pg.GetDB().Exec("INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)", email, passwordHash, role)
	return err
}

func (userRepo *PGUserRepo) SetPasswordHashByEmail(email, passwordHash string) error {
	_, err := userRepo.pg.GetDB().Exec("UPDATE users SET password_hash = $1 WHERE email = $2", passwordHash, email)
	return err
}

func (userRepo *PGUserRepo) DeleteUserByID(id int) error {
	_, err := userRepo.pg.GetDB().Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func NewPGUserRepo(pg *Postgres) model.UserRepo {
	return &PGUserRepo{
		pg:         pg,
		userMapper: NewPGUserMapper(),
	}
}
