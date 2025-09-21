package server

import (
	"context"
	"errors"
	"github.com/AndreySirin/-Effective-Mobile-/internal/storage"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	lg      *slog.Logger
	srv     *http.Server
	storage storage.SubscriptionStorage
}

func New(log *slog.Logger, addr string, stor storage.SubscriptionStorage) *Server {
	lg := log.With("module", "server")
	lg.Info("initializing server", "addr", addr)

	s := &Server{
		lg:      lg,
		storage: stor,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			lg.Info("registering API routes")
			r.Post("/subs", s.CreateSubs)
			r.Get("/subs/{id}", s.ReadSubs)
			r.Post("/subs/{id}", s.UpdateSubs)
			r.Delete("/subs/{id}", s.DeleteSubs)
			r.Get("/subs", s.ListSubs)
			r.Post("/cost", s.TotalCost)
		})
	})

	s.srv = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	lg.Info("server initialized successfully")
	return s
}

func (s *Server) Run() error {
	s.lg.Info("starting server", "addr", s.srv.Addr)

	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		s.lg.Info("server stopped gracefully")
		return nil
	}
	if err != nil {
		s.lg.Error("server encountered an unexpected error", "err", err)
		return err
	}

	return nil
}

func (s *Server) ShutDown() error {
	s.lg.Info("shutting down server", "timeout_sec", 3)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.lg.Error("server shutdown failed", "err", err)
		return err
	}

	s.lg.Info("server shutdown completed successfully")
	return nil
}
