package postgres

import (
	"context"
	"errors"
	"fmt"

	"socialite/database"

	"github.com/jackc/pgx/v5/pgconn"
)

func (c *Client) PutUser(ctx context.Context, user *database.User) error {
	if user == nil {
		return errors.New("user input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	_, err := c.Pool.Exec(
		queryCtx,
		`INSERT INTO users 
			(name, created_at) 
		VALUES 
			($1, $2)`,
		user.Name, user.CreatedAt,
	)
	if err != nil {
		// duplicate user name check
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return database.Err_DuplicatePrimaryKey
		}
		return fmt.Errorf("executing postgres insertion: %s", err.Error())
	}
	return nil
}

func (c *Client) GetUser(ctx context.Context, name string) (*database.User, error) {
	if name == "" {
		return nil, errors.New("name input is empty")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	row := c.Pool.QueryRow(
		queryCtx,
		`SELECT 
			created_at
		FROM users
		WHERE 
			name = $1`,
		name,
	)

	user := &database.User{
		Name: name,
	}
	err := row.Scan(&user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("scanning postgres row: %s", err.Error())
	}
	return user, nil
}
