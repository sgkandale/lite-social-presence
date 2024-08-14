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

func TestGetFriendshipById(t *testing.T) {
	friendship, err := dbConn.GetFriendshipById(
		context.Background(),
		801,
	)
	if err != nil {
		t.Error(err)
		return
	}
	log.Printf("friendship : %+v", friendship)
}

var testParty = &database.Party{
	Name:      "party1",
	Creator:   "user1",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func TestPutPartyOne(t *testing.T) {
	err := dbConn.PutParty(
		context.Background(),
		testParty,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGetParty(t *testing.T) {
	party, err := dbConn.GetParty(
		context.Background(),
		testParty.Name,
	)
	if err != nil {
		t.Error(err)
		return
	}
	log.Printf("party : %+v", party)
}

func TestGetCreatedParties(t *testing.T) {
	parties, err := dbConn.GetCreatedParties(
		context.Background(),
		testParty.Creator,
	)
	if err != nil {
		t.Error(err)
		return
	}
	if len(parties) == 0 {
		log.Print("parties is empty")
		return
	}
	for _, eachParty := range parties {
		log.Printf("party : %+v", eachParty)
	}
	log.Printf("total parties created : %d", len(parties))
}

var testPartyMembership = &database.PartyMembership{
	PartyName: "party1",
	UserName:  "user1",
	Status:    database.PartyMembership_Status_Invited,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func TestPutPartyMembership(t *testing.T) {
	err := dbConn.PutPartyMembership(
		context.Background(),
		testPartyMembership,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdatePartyMembership(t *testing.T) {
	testPartyMembership.Status = database.PartyMembership_Status_Active
	err := dbConn.UpdatePartyMembership(
		context.Background(),
		testPartyMembership,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestDeletePartyMembership(t *testing.T) {
	err := dbConn.DeletePartyMembership(
		context.Background(),
		testPartyMembership,
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPartyMembership(t *testing.T) {
	membership, err := dbConn.GetPartyMembership(
		context.Background(),
		testPartyMembership.PartyName,
		testPartyMembership.UserName,
	)
	if err != nil {
		t.Error(err)
		return
	}
	log.Printf("membership : %+v", membership)
}

func TestGetPendingFriendRequests(t *testing.T) {
	requests, err := dbConn.GetPendingFriendRequests(
		context.Background(),
		"user_2",
	)
	if err != nil {
		t.Error(err)
		return
	}
	if len(requests) == 0 {
		log.Print("requests is empty")
		return
	}
	for _, eachRequest := range requests {
		log.Printf("request : %+v", eachRequest)
	}
	log.Printf("total requests : %d", len(requests))
}

func TestGetUserFriendsList(t *testing.T) {
	friendsMap, err := dbConn.GetUserFriendsList(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if len(friendsMap) == 0 {
		log.Print("friendsMap is empty")
		return
	}
	for userName, friendsList := range friendsMap {
		log.Printf("%s -> %+v", userName, friendsList)
	}
}

func TestGetPartyMembers(t *testing.T) {
	members, err := dbConn.GetPartyMembers(
		context.Background(),
		"party_2",
	)
	if err != nil {
		t.Error(err)
		return
	}
	if len(members) == 0 {
		log.Print("members is empty")
		return
	}
	for _, eachMember := range members {
		log.Printf("member : %+v", eachMember)
	}
	log.Printf("total members : %d", len(members))
}

func TestGetAllPartyMembers(t *testing.T) {
	members, err := dbConn.GetAllPartyMembers(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if len(members) == 0 {
		log.Print("members is empty")
		return
	}
	for partyName, membersList := range members {
		log.Printf("%s -> %+v", partyName, membersList)
	}
	log.Printf("total parties : %d", len(members))
}
