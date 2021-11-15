package server

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/config"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:			":" + cfg.HTTP.Port,
			Handler:		handler,
			MaxHeaderBytes:	1 << 20,
			ReadTimeout:	10 * time.Second,
			WriteTimeout:	10 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
