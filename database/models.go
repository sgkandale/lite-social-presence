package database

import (
	"errors"
	"time"
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
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}
