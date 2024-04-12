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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS migrations (index INT)")

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

	err = pg.migrate(tx, migrationsDir, runTestMigrations)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (pg *Postgres) migrate(tx *sqlx.Tx, migrationsDir string, runTestMigrations bool) error {
	migrationsPattern := migrationsDir + "/*.sql"
	migrations, _ := filepath.Glob(migrationsPattern)

	// Subtract test migrations
	if !runTestMigrations {
		testMigrationsPattern := migrationsDir + "/test.*.sql"
		testMigrations, _ := filepath.Glob(testMigrationsPattern)
		migrations = util.Diff(migrations, testMigrations)
	}

	lastMigIndex := pg.getLastMigrationIndex()

	if lastMigIndex+1 > len(migrations) {
		return fmt.Errorf("found %d migrations on disk, but DB claims to have run %d migrations. Some migrations must have been deleted", len(migrations), lastMigIndex+1)
	}

	remainingMigrations := migrations[lastMigIndex+1:]

	for i, path := range remainingMigrations {
		fmt.Printf("Running migration %s\n", path)

		normalizedIndex := i + lastMigIndex + 1
		contentBytes, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		content := string(contentBytes)

		_, err = tx.Exec(content)

		if err != nil {
			return err
		}

		err = pg.insertMigration(tx, normalizedIndex)
		if err != nil {
			return err
		}
	}

	count, err := pg.getNumMigrations(tx)

	if err != nil {
		return err
	}

	if count != len(migrations) {
		return fmt.Errorf("found %d migrations on disk, but DB claims to have run %d migrations. Some migrations must have been deleted", len(migrations), count)
	}

	return nil
}

func (pg *Postgres) insertMigration(tx *sqlx.Tx, i int) error {
	_, err := tx.Exec("INSERT INTO migrations (index) VALUES ($1)", i)
	return err
}

func (pg *Postgres) getNumMigrations(tx *sqlx.Tx) (int, error) {
	var count []int
	err := tx.Select(&count, "SELECT COUNT(*) FROM migrations")

	if err != nil || len(count) == 0 {
		return 0, err
	}

	return count[0], nil
}

func (pg *Postgres) getLastMigrationIndex() int {
	index := []int{}
	err := pg.db.Select(&index, "SELECT index FROM migrations ORDER BY index DESC LIMIT 1")

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
