package database

import (
	"errors"
	"strings"
	"time"
)

type Friendship_Status string

const (
	Friendship_Status_Sent      Friendship_Status = "sent"
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
	// sort user names in lexicographical order
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	return &Friendship{
		User1:     strings.ToLower(user1),
		User2:     strings.ToLower(user2),
		Status:    Friendship_Status_Sent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
