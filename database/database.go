package database

import "context"

type Database interface {
	// user methods
	PutUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, name string) (*User, error)

	// friends methods
	GetUserFriends(ctx context.Context, name string) ([]*User, error)
	PutFriendship(ctx context.Context, friendship *Friendship) error
	GetFriendship(ctx context.Context, user1, user2 string) (*Friendship, error)
	GetFriendshipById(ctx context.Context, friendshipId int32) (*Friendship, error)
	UpdateFriendship(ctx context.Context, friendship *Friendship) error
	DeleteFriendship(ctx context.Context, friendshipId int32) error
}
