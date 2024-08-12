package server

import (
	"context"
	"fmt"
	"log"

	"socialite/config"
	"socialite/database"
	"socialite/database/postgres"

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

	// connections
	db database.Database
}

func New(ctx context.Context, cfg *config.Config) *Server {
	if cfg == nil {
		log.Fatal("[ERROR] server.New: config is nil")
	}

	ginEngine := gin.Default()

	var dbCnn database.Database
	if cfg.Database.Type == "postgres" {
		dbCnn = postgres.New(ctx, &cfg.Database)
	} else {
		log.Fatal("[ERROR] database type is not supported: ", cfg.Database.Type)
	}

	return &Server{
		engine:      ginEngine,
		port:        cfg.Server.Port,
		name:        cfg.Server.ServiceName,
		tls:         cfg.Server.TLS,
		tlsCertPath: cfg.Server.CertPath,
		tlsKeyPath:  cfg.Server.KeyPath,
		db:          dbCnn,
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
