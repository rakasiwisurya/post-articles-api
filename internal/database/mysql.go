package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"post-articles-api/internal/config"
)

// EnsureDatabase creates the application database when it does not exist yet,
// so a fresh environment only needs a running MySQL server.
func EnsureDatabase(cfg config.Config) error {
	db, err := sql.Open("mysql", cfg.ServerDSN())
	if err != nil {
		return fmt.Errorf("open mysql server connection: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci",
		cfg.DBName,
	)
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("create database %s: %w", cfg.DBName, err)
	}
	return nil
}

// Connect opens a pooled connection to the application database and verifies
// it is reachable.
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return db, nil
}
