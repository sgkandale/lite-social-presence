package database

import "context"

type Database interface {
	// user methods
	PutUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, name string) (*User, error)
}
