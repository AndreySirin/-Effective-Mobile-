package main

import (
	"context"
	"github.com/AndreySirin/-Effective-Mobile-/internal/config"
	"github.com/AndreySirin/-Effective-Mobile-/internal/logger"
	"github.com/AndreySirin/-Effective-Mobile-/internal/server"
	"github.com/AndreySirin/-Effective-Mobile-/internal/storage"
	migrate "github.com/rubenv/sql-migrate"
	"os"
	"os/signal"
	"sync"
)

func main() {
	lg := logger.New()
	lg.Info("starting application")
	cfg, err := config.Load(lg)
	if err != nil {
		lg.Error("error loading config", "error", err)
		return
	}
	lg.Info("config loaded",
		"db_host", cfg.Postgres.Address,
		"db_name", cfg.Postgres.DbName,
		"http_port", cfg.Server.Port,
	)
	db, err := storage.New(lg,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Address,
		cfg.Postgres.DbName)
	if err != nil {
		lg.Error("error connecting to database", "error", err)
	}
	lg.Info("database connection established")

	defer func() {
		err = db.Close()
		if err != nil {
			lg.Error("error closing database connection", "error", err)
		} else {
			lg.Info("database connection closed")
		}
	}()
	err = db.Migrate(migrate.Up)
	if err != nil {
		lg.Error("error migrating database", "error", err)
	}
	lg.Info("database migration complete")
	srv := server.New(lg, cfg.Server.Port, db)
	lg.Info("server initialized", "port", cfg.Server.Port)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		lg.Info("server starting...")
		err = srv.Run()
		if err != nil {
			lg.Error("error running server", "error", err)
			return
		}
		lg.Info("server stopped gracefully")
		wg.Done()
	}()

	go func() {
		<-ctx.Done()
		lg.Info("shutdown signal received")
		err = srv.ShutDown()
		if err != nil {
			lg.Error("error shutting down server", "error", err)
		} else {
			lg.Info("server shutdown gracefully")
		}
		wg.Done()
	}()
	wg.Wait()
	lg.Info("server shutdown complete")
}
