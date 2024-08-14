package cache

import (
	"context"
	"time"
)

type Cache interface {
	PutUserOnline(ctx context.Context, userName string) error
	IsUserOnline(ctx context.Context, userName string) (bool, error)

	PutUserFriendsList(ctx context.Context, userName string, friends []string) error
	GetUserFriendsList(ctx context.Context, userName string) ([]string, error)

	PutPartyMembersList(ctx context.Context, partyName string, members []string) error
	GetPartyMembersList(ctx context.Context, partyName string) ([]string, error)
}

const (
	UserOnlineExpiry = time.Second * 10
)
