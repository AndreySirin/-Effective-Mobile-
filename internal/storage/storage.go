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
	lg = lg.With("module", "storage")
	lg.Info("initializing database connection", "user", username, "host", address, "db", database)

	dsn := (&url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(username, password),
		Host:   address,
		Path:   database,
	}).String()

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		lg.Error("failed to open database connection", "err", err, "dsn", dsn)
		return nil, fmt.Errorf("init db: %v", err)
	}
	lg.Info("database connection opened successfully")

	if err = sqlDB.Ping(); err != nil {
		lg.Error("failed to ping database", "err", err)
		return nil, fmt.Errorf("ping db: %v", err)
	}
	lg.Info("database ping successful")

	return &Storage{
		lg: lg,
		db: sqlDB,
	}, nil
}

func (s *Storage) Close() error {
	s.lg.Info("closing database connection")
	err := s.db.Close()
	if err != nil {
		s.lg.Error("error closing database connection", "err", err)
		return err
	}
	s.lg.Info("database connection closed successfully")
	return nil
}

func (s *Storage) Migrate(direction migrate.MigrationDirection) error {
	s.lg.Info("starting database migration", "direction", direction)

	migrations := &migrate.FileMigrationSource{
		Dir: "/root/migrate",
	}

	n, err := migrate.Exec(s.db, "postgres", migrations, direction)
	if err != nil {
		s.lg.Error("database migration failed", "err", err)
		return fmt.Errorf("error for migrate: %v", err)
	}

	s.lg.Info("database migration completed successfully", "migrations_applied", n)
	return nil
}
