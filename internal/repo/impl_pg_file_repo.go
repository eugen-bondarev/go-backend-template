package repo

import (
	"errors"
	"go-backend-template/internal/postgres"
)

type PGFileRepo struct {
	pg *postgres.Postgres
}

func (fileRepo *PGFileRepo) GetPath(ID FileHandle) (string, error) {
	var matches []string

	err := fileRepo.pg.GetDB().Select(&matches, "SELECT path FROM files WHERE id = $1", ID)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", errors.New("file not found")
	}

	return matches[0], nil
}

func (fileRepo *PGFileRepo) AddPath(path string) (FileHandle, error) {
	matches := make([]FileHandle, 0, 1)
	err := fileRepo.pg.GetDB().Select(&matches, "INSERT INTO files (path) VALUES ($1) RETURNING id", path)
	if err != nil {
		return -1, err
	}

	if len(matches) == 0 {
		return -1, errors.New("failed to retrieve the ID of a new item")
	}

	return matches[0], nil
}

func NewPGFileRepo(pg *postgres.Postgres) IFileRepo {
	return &PGFileRepo{
		pg: pg,
	}
}
