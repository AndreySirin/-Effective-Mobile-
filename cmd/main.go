package main

import (
	"context"
	"github.com/AndreySirin/-Effective-Mobile-/internal/config"
	"github.com/AndreySirin/-Effective-Mobile-/internal/logger"
	"github.com/AndreySirin/-Effective-Mobile-/internal/server"
	"github.com/AndreySirin/-Effective-Mobile-/internal/storage"
	"os"
	"os/signal"
	"sync"
)

func main() {
	lg := logger.New()
	cfg, err := config.Load()
	if err != nil {
		lg.Error("error loading config:", err)
	}
	db, err := storage.New(lg,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Address,
		cfg.Postgres.DbName)
	if err != nil {
		lg.Error("error connecting to database:", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			lg.Error("error closing database connection:", err)
		}
	}()
	srv := server.New(lg, cfg.Server.Port)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		err = srv.Run()
		if err != nil {
			lg.Error("error running server:", err)
			return
		}
		wg.Done()
	}()

	go func() {
		<-ctx.Done()
		err = srv.ShutDown()
		if err != nil {
			lg.Error("error shutting down server:", err)
		}
		wg.Done()
	}()
	wg.Wait()

}
