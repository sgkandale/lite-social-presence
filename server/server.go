package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"socialite/cache"
	"socialite/cache/state"
	"socialite/config"
	"socialite/database"
	"socialite/database/postgres"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	db       database.Database
	cache    cache.Cache
	upgrader websocket.Upgrader

	// internal variables
	rwmutex               sync.RWMutex
	userOnlineStatus      chan string
	userWebsocketChannels map[string]chan []byte
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

	var cacheConn cache.Cache
	if cfg.Cache.Type == "state" {
		cacheConn = state.New(ctx, &cfg.Cache)
	} else {
		log.Fatal("[ERROR] cache type is not supported: ", cfg.Cache.Type)
	}

	return &Server{
		engine:      ginEngine,
		port:        cfg.Server.Port,
		name:        cfg.Server.ServiceName,
		tls:         cfg.Server.TLS,
		tlsCertPath: cfg.Server.CertPath,
		tlsKeyPath:  cfg.Server.KeyPath,
		db:          dbCnn,
		cache:       cacheConn,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		rwmutex:               sync.RWMutex{},
		userOnlineStatus:      make(chan string, 1_000),
		userWebsocketChannels: make(map[string]chan []byte, 1_000),
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
