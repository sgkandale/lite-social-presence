package postgres

import (
	"context"
	"log"
	"strings"
	"time"

	"socialite/config"
	"socialite/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	*pgxpool.Pool
	timeout time.Duration
}

func New(ctx context.Context, cfg *config.DatabaseConfig) database.Database {
	if cfg == nil {
		log.Fatal("[ERROR] postgres.New: config is nil")
	}
	if !strings.EqualFold(cfg.Type, "postgres") {
		log.Fatal("[ERROR] postgres.New: invalid database type:", cfg.Type)
	}

	timeout := time.Second * time.Duration(cfg.Timeout)
	connectCtx, cancelConnectCtx := context.WithTimeout(ctx, timeout)
	defer cancelConnectCtx()

	pgPool, err := pgxpool.New(connectCtx, cfg.UriString)
	if err != nil {
		log.Fatal("[ERROR] postgres.New: creating client:", err)
	}

	pingCtx, cancelPingCtx := context.WithTimeout(ctx, timeout)
	defer cancelPingCtx()

	err = pgPool.Ping(pingCtx)
	if err != nil {
		log.Fatal("[ERROR] postgres.New: pinging db:", err)
	}

	return &Client{
		Pool:    pgPool,
		timeout: timeout,
	}
}
