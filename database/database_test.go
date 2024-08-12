package database_test

import (
	"context"
	"log"
	"testing"
	"time"

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

func TestPutUser(t *testing.T) {
	err := dbConn.PutUser(
		context.Background(),
		&database.User{
			Name:      "user1",
			CreatedAt: time.Now(),
		},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUser(t *testing.T) {
	user, err := dbConn.GetUser(
		context.Background(),
		"user1",
	)
	if err != nil {
		t.Error(err)
		return
	}
	if user == nil {
		t.Error("user is nil")
		return
	}
	if user.Name != "user1" {
		t.Error("user name is incorrect")
	}
	if user.CreatedAt.IsZero() {
		t.Error("user created at is zero")
	}
	log.Printf("user : %+v", user)
}
