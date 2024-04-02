package postgres

import (
	"context"
	"fmt"
	"go-backend-template/internal/util"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(host, username, password, port string) (Postgres, error) {
	connectionStr := fmt.Sprintf(
		"user=%s host=%s port=%s password=%s sslmode=disable",
		username,
		host,
		port,
		password,
	)
	return NewPostgresFromConnectionString(connectionStr)
}

func NewPostgresFromConnectionString(connectionString string) (Postgres, error) {
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		return Postgres{}, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS migrations (index INT, hash VARCHAR)")

	if err != nil {
		return Postgres{}, err
	}

	return Postgres{
		db,
	}, nil
}

func (pg *Postgres) getHashFromMigrationsDBAt(i int) (string, bool) {
	var hashFromDB []string
	err := pg.db.Select(&hashFromDB, "SELECT hash FROM migrations WHERE index = $1", i)

	if len(hashFromDB) == 0 || err != nil {
		return "", false
	}

	return hashFromDB[0], true
}

func (pg *Postgres) insertHash(i int, hash string) error {
	_, err := pg.db.Exec("INSERT INTO migrations (index, hash) VALUES ($1, $2)", i, hash)
	return err
}

func (pg *Postgres) getNumMigrations() (int, error) {
	var count []int
	err := pg.db.Select(&count, "SELECT COUNT(*) FROM migrations")

	if err != nil || len(count) == 0 {
		return 0, err
	}

	return count[0], nil
}

func (pg *Postgres) Migrate(migrationsDir string) error {
	ctx := context.Background()

	tx, err := pg.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	migrationsPattern := migrationsDir + "/*.sql"
	migrations, err := filepath.Glob(migrationsPattern)

	if err != nil {
		return err
	}

	for i, path := range migrations {
		fmt.Printf("Running migration %s\n", path)
		contentBytes, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		content := string(contentBytes)

		hashNum := util.HashFNV32a(content)
		hash := fmt.Sprintf("%v", hashNum)
		hashFromDB, exists := pg.getHashFromMigrationsDBAt(i)

		if exists && hashFromDB != hash {
			panic("migration " + path + " was updated, which is not allowed")
		}

		_, err = pg.db.Exec(content)

		if err != nil {
			return err
		}

		if exists {
			continue
		}

		err = pg.insertHash(i, hash)

		if err != nil {
			return err
		}
	}

	count, err := pg.getNumMigrations()

	if err != nil {
		return err
	}

	if count != len(migrations) {
		panic(fmt.Sprintf("found %d migrations on disk, but DB claims to have run %d migrations. Some migrations must have been deleted", len(migrations), count))
	}

	return tx.Commit()
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
