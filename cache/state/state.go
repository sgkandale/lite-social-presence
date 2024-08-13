package state

import (
	"context"
	"log"
	"strings"

	"socialite/cache"
	"socialite/config"

	"github.com/dgraph-io/ristretto"
)

type Client struct {
	cache *ristretto.Cache
}

func New(ctx context.Context, cfg *config.CacheConfig) cache.Cache {
	if cfg == nil {
		log.Fatal("[ERROR] cache config is nil")
	}
	if !strings.EqualFold(cfg.Type, "state") {
		log.Fatal("[ERROR] cache type is unknown")
	}

	cache, err := ristretto.NewCache(
		&ristretto.Config{
			NumCounters: 1e5,
			MaxCost:     104_857_600,
			BufferItems: 64,
		},
	)
	if err != nil {
		log.Fatal("[ERROR] creating ristretto cache : ", err.Error())
	}
	return &Client{
		cache: cache,
	}
}
