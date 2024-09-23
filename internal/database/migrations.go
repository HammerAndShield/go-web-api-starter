package database

import (
	"embed"
	"github.com/pressly/goose/v3"
)

const (
	DialectPostgres = "postgres"
	DialectMySQL    = "mysql"
	DialectSQLite   = "sqlite3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (db *DB) RunMigrations(dialect string) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
