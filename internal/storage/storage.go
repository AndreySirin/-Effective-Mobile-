package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"log/slog"
	"net/url"
)

type Storage struct {
	lg *slog.Logger
	db *sql.DB
}

func New(
	lg *slog.Logger,
	username string,
	password string,
	address string,
	database string,
) (*Storage, error) {
	dsn := (&url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(username, password),
		Host:   address,
		Path:   database,
	}).String()

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("init db: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %v", err)
	}

	return &Storage{
		lg: lg.With("module", "storage"),
		db: sqlDB,
	}, nil
}
func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Migrate(bool migrate.MigrationDirection) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "/root/migrate",
	}
	_, err := migrate.Exec(s.db, "postgres", migrations, bool)
	if err != nil {
		return fmt.Errorf("error for migrate: %v", err)
	}
	return nil
}
