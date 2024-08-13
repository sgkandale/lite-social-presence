package cache_test

import (
	"context"
	"testing"

	"socialite/cache"
	"socialite/cache/state"
	"socialite/config"
)

var cacheConn cache.Cache

func init() {
	cacheConn = state.New(
		context.Background(),
		&config.CacheConfig{
			Type: "state",
		},
	)
}

func TestNew(t *testing.T) {
	if cacheConn == nil {
		t.Error("cacheConn is nil")
	}
}
