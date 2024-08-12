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

func TestGetUserFriends(t *testing.T) {
	username := "user1"
	friendsList, err := dbConn.GetUserFriends(
		context.Background(),
		username,
	)
	if err != nil {
		t.Error(err)
		return
	}
	if friendsList == nil {
		return
	}
	if len(friendsList) == 0 {
		log.Print("friendsList is empty")
		return
	}
	for _, eachFriend := range friendsList {
		log.Printf("friend : %+v", eachFriend)
	}
	log.Printf("total friends for user %s : %d", username, len(friendsList))
}

var testFriendship = &database.Friendship{
	User1:     "user1",
	User2:     "user2",
	Status:    database.Friendship_Status_Sent,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func TestPutFriendship(t *testing.T) {
	err := dbConn.PutFriendship(
		context.Background(),
		testFriendship,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateFriendship(t *testing.T) {
	testFriendship.Status = database.Friendship_Status_Confirmed
	err := dbConn.UpdateFriendship(
		context.Background(),
		testFriendship,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteFriendship(t *testing.T) {
	err := dbConn.DeleteFriendship(
		context.Background(),
		701,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFriendship(t *testing.T) {
	friendship, err := dbConn.GetFriendship(
		context.Background(),
		"user1",
		"user2",
	)
	if err != nil {
		t.Error(err)
		return
	}
	log.Printf("friendship : %+v", friendship)
}
