package postgres

import (
	"context"

	"socialite/database"
)

func (c *Client) PutUser(ctx context.Context, user *database.User) error {
	return nil
}
func (c *Client) GetUser(ctx context.Context, name string) (*database.User, error) {
	return nil, nil
}
