package server

import (
	"context"
	"net/http"
	"time"

	"example.com/go-crud/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, router *gin.Engine) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.ServerAddress,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
