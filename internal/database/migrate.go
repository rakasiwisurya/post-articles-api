package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"post-articles-api/migrations"
)

// Migrate applies all pending migrations from the embedded migration files.
// It is a no-op when the schema is already up to date.
func Migrate(db *sql.DB, dbName string) error {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("load migration source: %w", err)
	}

	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{DatabaseName: dbName})
	if err != nil {
		return fmt.Errorf("init migration driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
