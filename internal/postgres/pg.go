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

func createMigrationsTable(db *sqlx.DB, migrationsTable string) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + migrationsTable + " (index INT)")
	return err
}

func NewPostgresFromConnectionString(connectionString string) (Postgres, error) {
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		return Postgres{}, err
	}

	return Postgres{
		db,
	}, nil
}

func (pg *Postgres) Migrate(migrationsDir string, runTestMigrations bool) error {
	ctx := context.Background()
	tx, err := pg.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	err = pg.migrate(pg.db, tx, migrationsDir, runTestMigrations)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func getMigrations(dir string, pattern string) []string {
	migrationsPattern := dir + "/" + pattern
	migrations, _ := filepath.Glob(migrationsPattern)
	return migrations
}

func (pg *Postgres) migrate(db *sqlx.DB, tx *sqlx.Tx, migrationsDir string, runTestMigrations bool) error {
	testMigrations := getMigrations(migrationsDir, "test.*.sql")
	migrations := util.Diff(getMigrations(migrationsDir, "*.sql"), testMigrations)

	return util.EvalUntilErr([]func() error{
		func() error {
			return createMigrationsTable(db, "migrations")
		},
		func() error {
			return pg.migrateFiles(tx, migrations, "migrations")
		},
		func() error {
			if !runTestMigrations {
				return nil
			}
			return createMigrationsTable(db, "test_migrations")
		},
		func() error {
			if !runTestMigrations {
				return nil
			}
			return pg.migrateFiles(tx, testMigrations, "test_migrations")
		},
	})
}

func (pg *Postgres) migrateFiles(tx *sqlx.Tx, files []string, migrationsTable string) error {
	lastMigIndex := pg.getLastMigrationIndex(migrationsTable)

	if lastMigIndex+1 == len(files) {
		return nil
	}

	if lastMigIndex+1 > len(files) {
		return fmt.Errorf("found %d migrations on disk, but DB claims to have run %d migrations. Some migrations must have been deleted", len(files), lastMigIndex+1)
	}

	for i, path := range files[lastMigIndex+1:] {
		fmt.Printf("Running migration %s\n", path)

		contentBytes, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		content := string(contentBytes)

		_, err = tx.Exec(content)

		if err != nil {
			return err
		}

		err = pg.insertMigration(tx, migrationsTable, i+lastMigIndex+1)
		if err != nil {
			return err
		}
	}

	count, err := pg.getNumMigrations(tx, migrationsTable)

	if err != nil {
		return err
	}

	if count != len(files) {
		return fmt.Errorf("found %d migrations on disk, but DB claims to have run %d migrations. Some migrations must have been deleted", len(files), count)
	}

	return nil
}

func (pg *Postgres) insertMigration(tx *sqlx.Tx, migrationsTable string, i int) error {
	_, err := tx.Exec("INSERT INTO "+migrationsTable+" (index) VALUES ($1)", i)
	return err
}

func (pg *Postgres) getNumMigrations(tx *sqlx.Tx, migrationsTable string) (int, error) {
	var count []int
	err := tx.Select(&count, "SELECT COUNT(*) FROM "+migrationsTable)

	if err != nil || len(count) == 0 {
		return 0, err
	}

	return count[0], nil
}

func (pg *Postgres) getLastMigrationIndex(migrationsTable string) int {
	index := []int{}
	err := pg.db.Select(&index, "SELECT index FROM "+migrationsTable+" ORDER BY index DESC LIMIT 1")

	if err != nil {
		return -1
	}

	if len(index) == 0 {
		return -1
	}

	return index[0]
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
