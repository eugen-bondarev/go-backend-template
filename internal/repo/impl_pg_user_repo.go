package repo

import (
	"errors"
	"go-backend-template/internal/model"
	"go-backend-template/internal/postgres"
)

type PGUserRepo struct {
	pg         *postgres.Postgres
	userMapper model.ModelMapper[model.User, postgres.PGUser]
}

func (userRepo *PGUserRepo) getUserByEmail(email string) (postgres.PGUser, error) {
	var users []postgres.PGUser

	err := userRepo.pg.GetDB().Select(&users, "SELECT * FROM users WHERE email = $1", email)

	if err != nil {
		return postgres.PGUser{}, err
	}

	if len(users) == 0 {
		return postgres.PGUser{}, errors.New("user not found")
	}

	return users[0], nil
}

func (userRepo *PGUserRepo) getUserByID(ID int) (postgres.PGUser, error) {
	var users []postgres.PGUser

	err := userRepo.pg.GetDB().Select(&users, "SELECT * FROM users WHERE id = $1", ID)

	if err != nil {
		return postgres.PGUser{}, err
	}

	if len(users) == 0 {
		return postgres.PGUser{}, errors.New("user not found")
	}

	return users[0], nil
}

func (userRepo *PGUserRepo) getUsers() ([]postgres.PGUser, error) {
	var users []postgres.PGUser

	err := userRepo.pg.GetDB().Select(&users, "SELECT * FROM users")

	if err != nil {
		return []postgres.PGUser{}, err
	}

	if len(users) == 0 {
		return []postgres.PGUser{}, errors.New("user not found")
	}

	return users, nil
}

func (userRepo *PGUserRepo) GetUserByEmail(email string) (model.User, error) {
	user, err := userRepo.getUserByEmail(email)

	if err != nil {
		return model.User{}, err
	}

	return userRepo.userMapper.ToModel(user), nil
}

func (userRepo *PGUserRepo) GetUserByID(ID int) (model.User, error) {
	user, err := userRepo.getUserByID(ID)

	if err != nil {
		return model.User{}, err
	}

	return userRepo.userMapper.ToModel(user), nil
}

func (userRepo *PGUserRepo) GetUsers() ([]model.User, error) {
	users, err := userRepo.getUsers()

	if err != nil {
		return []model.User{}, err
	}

	return model.ManyToModel(
		userRepo.userMapper,
		users,
	), nil
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

func NewPGUserRepo(pg *postgres.Postgres) IUserRepo {
	return &PGUserRepo{
		pg:         pg,
		userMapper: postgres.NewPGUserMapper(),
	}
}
