package server

import (
	"context"
	"fmt"
	"log"

	"socialite/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	// server internals
	engine      *gin.Engine
	port        int
	name        string
	tls         bool
	tlsCertPath string
	tlsKeyPath  string
}

func New(ctx context.Context, cfg *config.Config) *Server {
	ginEngine := gin.Default()

	return &Server{
		engine:      ginEngine,
		port:        cfg.Server.Port,
		name:        cfg.Server.ServiceName,
		tls:         cfg.Server.TLS,
		tlsCertPath: cfg.Server.CertPath,
		tlsKeyPath:  cfg.Server.KeyPath,
	}
}

func (s *Server) Start() error {
	log.Printf("[INFO] starting server for %s on port %d with tls %t", s.name, s.port, s.tls)
	// start the server
	port := fmt.Sprintf(":%d", s.port)
	if s.tls {
		return s.engine.RunTLS(port, s.tlsCertPath, s.tlsKeyPath)
	}
	return s.engine.Run(port)
}
