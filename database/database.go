package database

import "context"

type Database interface {
	// user methods
	PutUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, name string) (*User, error)

	// friends methods
	GetUserFriends(ctx context.Context, name string) ([]*User, error)
	PutFriendship(ctx context.Context, friendship *Friendship) error
	GetPendingFriendRequests(ctx context.Context, userName string) ([]*Friendship, error)
	GetFriendship(ctx context.Context, user1, user2 string) (*Friendship, error)
	GetFriendshipById(ctx context.Context, friendshipId int32) (*Friendship, error)
	UpdateFriendship(ctx context.Context, friendship *Friendship) error
	DeleteFriendship(ctx context.Context, friendshipId int32) error

	// party methods
	PutParty(ctx context.Context, party *Party) error
	GetParty(ctx context.Context, partyName string) (*Party, error)
	GetCreatedParties(ctx context.Context, userName string) ([]*Party, error)

	// party membership methods
	PutPartyMembership(ctx context.Context, membership *PartyMembership) error
	GetPartyMembership(ctx context.Context, partyName, userName string) (*PartyMembership, error)
	UpdatePartyMembership(ctx context.Context, membership *PartyMembership) error
	DeletePartyMembership(ctx context.Context, membership *PartyMembership) error
}
