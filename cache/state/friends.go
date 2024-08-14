package state

import (
	"context"
)

func (c *Client) PutUserFriendsList(ctx context.Context, userName string, friends []string) error {
	c.cache.SetWithTTL(
		UserFriendsListKey(userName),
		friends,
		1,
		0,
	)
	return nil
}

func (c *Client) GetUserFriendsList(ctx context.Context, userName string) ([]string, error) {
	// get from cache
	if friends, ok := c.cache.Get(UserFriendsListKey(userName)); ok {
		return friends.([]string), nil
	}
	return nil, nil
}
