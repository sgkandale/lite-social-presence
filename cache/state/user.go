package state

import (
	"context"

	"socialite/cache"
)

func (c *Client) PutUserOnline(ctx context.Context, userName string) error {
	c.cache.SetWithTTL(
		userName,
		true,
		1,
		cache.UserOnlineExpiry,
	)
	return nil
}

func (c *Client) IsUserOnline(ctx context.Context, userName string) (bool, error) {
	userOnline, ok := c.cache.Get(userName)
	if ok {
		return userOnline.(bool), nil
	}
	return false, nil
}
