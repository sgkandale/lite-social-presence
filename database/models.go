package database

import (
	"errors"
	"strings"
	"time"
)

type Friendship_Status string

const (
	Friendship_Status_Sent      Friendship_Status = "sent"
	Friendship_Status_Received  Friendship_Status = "received"
	Friendship_Status_Confirmed Friendship_Status = "confirmed"
)

type User struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("user name is empty")
	}
	return &User{
		Name:      strings.ToLower(name),
		CreatedAt: time.Now(),
	}, nil
}

type Friendship struct {
	Id        int32             `json:"id"`
	User1     string            `json:"user1"`
	User2     string            `json:"user2"`
	Status    Friendship_Status `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func NewFriendship(user1, user2 string) (*Friendship, error) {
	if user1 == "" || user2 == "" {
		return nil, errors.New("user name is empty")
	}
	return &Friendship{
		User1:     strings.ToLower(user1),
		User2:     strings.ToLower(user2),
		Status:    Friendship_Status_Sent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

type Party struct {
	Name      string    `json:"name"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewParty(name, creator string) (*Party, error) {
	if name == "" {
		return nil, errors.New("party name is empty")
	}
	if creator == "" {
		return nil, errors.New("creator name is empty")
	}

	return &Party{
		Name:      strings.ToLower(name),
		Creator:   strings.ToLower(name),
		CreatedAt: time.Now(),
	}, nil
}

type PartyMembership_Status string

const (
	PartyMembership_Status_Invited PartyMembership_Status = "invited"
	PartyMembership_Status_Active  PartyMembership_Status = "active"
)

type PartyMembership struct {
	PartyName string                 `json:"party_name"`
	UserName  string                 `json:"user_name"`
	Status    PartyMembership_Status `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func NewPartyMembership(partyName, userName string) (*PartyMembership, error) {
	if partyName == "" {
		return nil, errors.New("party name is empty")
	}
	if userName == "" {
		return nil, errors.New("user name is empty")
	}
	return &PartyMembership{
		PartyName: strings.ToLower(partyName),
		UserName:  strings.ToLower(userName),
		Status:    PartyMembership_Status_Invited,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
