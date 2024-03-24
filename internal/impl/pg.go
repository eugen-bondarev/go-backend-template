package impl

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(connectionString string) (Postgres, error) {
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		return Postgres{}, err
	}

	return Postgres{
		db,
	}, nil
}

func (pg *Postgres) ExecFile(path string) error {
	initRequestsBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	initRequestsSql := string(initRequestsBytes)
	_, err = pg.db.Exec(initRequestsSql)
	return err
}

func (pg *Postgres) GetDB() *sqlx.DB {
	return pg.db
}
