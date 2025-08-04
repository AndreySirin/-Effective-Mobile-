package server

import (
	"context"
	"errors"
	"fmt"
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
	s := &Server{
		lg:      log.With("module", "server"),
		storage: stor,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
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
	return s
}

func (s *Server) Run() error {
	s.lg.Info(fmt.Sprintf("Listening on %s", s.srv.Addr))
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		s.lg.Info("stopping the server")
	}
	return nil
}

func (s *Server) ShutDown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.lg.Error("error when stopping the server")
		return err
	}
	return nil
}
