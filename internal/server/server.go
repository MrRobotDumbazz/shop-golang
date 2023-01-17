package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	srv http.Server
}

func (s *Server) Start(port string, handlers chi.Router) error {
	s.srv = http.Server{
		Addr:         port,
		Handler:      handlers,
		WriteTimeout: time.Second * 3,
		ReadTimeout:  time.Second * 3,
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
