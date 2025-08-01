package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	lg  *slog.Logger
	srv *http.Server
}

func New(log *slog.Logger, addr string) *Server {
	s := &Server{
		lg: log.With("module", "server"),
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

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
