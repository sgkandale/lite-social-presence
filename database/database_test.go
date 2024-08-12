package database_test

import (
	"context"
	"testing"

	"socialite/config"
	"socialite/database"
	"socialite/database/postgres"
)

var dbConn database.Database

func init() {
	// create postgres client
	dbConn = postgres.New(
		context.Background(),
		&config.DatabaseConfig{
			Type:      "postgres",
			UriString: "",
			Timeout:   60,
		},
	)
}

func TestNew(t *testing.T) {
	if dbConn == nil {
		t.Error("dbConn is nil")
	}
}
