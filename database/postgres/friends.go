package postgres

import (
	"context"
	"errors"
	"fmt"

	"socialite/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (c *Client) GetUserFriends(ctx context.Context, name string) ([]*database.User, error) {
	if name == "" {
		return nil, errors.New("name input is empty")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	rows, err := c.Pool.Query(
		queryCtx,
		`SELECT
			user1, user2
		FROM friendships
		WHERE 
			status = $1
			AND (
				user1 = $2
				OR user2 = $2
			)`,
		database.Friendship_Status_Confirmed,
		name,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("querying postgres: %s", err.Error())
	}

	respUsers := make([]*database.User, 0)

	for rows.Next() {
		var user1, user2 string
		err = rows.Scan(&user1, &user2)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %s", err.Error())
		}

		otherUser := user1
		if otherUser == name {
			otherUser = user2
		}

		newUser, err := database.NewUser(otherUser)
		if err != nil {
			return nil, fmt.Errorf("creating new user while scanning each row: %s", err.Error())
		}
		respUsers = append(respUsers, newUser)
	}
	return respUsers, nil
}

func (c *Client) PutFriendship(ctx context.Context, friendship *database.Friendship) error {
	if friendship == nil {
		return errors.New("friendship input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	_, err := c.Pool.Exec(
		queryCtx,
		`INSERT INTO friendships 
			(user1, user2, status, created_at, updated_at)
		VALUES 
			($1, $2, $3, $4, $5)`,
		friendship.User1,
		friendship.User2,
		friendship.Status,
		friendship.CreatedAt,
		friendship.UpdatedAt,
	)
	if err != nil {
		// duplicate entry check
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return database.Err_DuplicatePrimaryKey
		}
		return fmt.Errorf("inserting friendship: %s", err.Error())
	}
	return nil
}

func (c *Client) UpdateFriendship(ctx context.Context, friendship *database.Friendship) error {
	if friendship == nil {
		return errors.New("friendship input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	_, err := c.Pool.Exec(
		queryCtx,
		`UPDATE friendships
		SET
			status = $1,
			updated_at = $2
		WHERE
			user1 = $3
			AND user2 = $4`,
		friendship.Status,
		friendship.UpdatedAt,
		friendship.User1,
		friendship.User2,
	)
	if err != nil {
		return fmt.Errorf("updating friendship: %s", err.Error())
	}

	return nil
}

func (c *Client) DeleteFriendship(ctx context.Context, friendshipId int32) error {
	if friendshipId == 0 {
		return errors.New("friendshipId input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	pgTag, err := c.Pool.Exec(
		queryCtx,
		`DELETE FROM friendships
		WHERE id = $1`,
		friendshipId,
	)
	if err != nil {
		return fmt.Errorf("deleting friendship: %s", err.Error())
	}
	if pgTag.RowsAffected() == 0 {
		return database.Err_NotFound
	}
	return nil
}

func (c *Client) GetFriendship(ctx context.Context, user1, user2 string) (*database.Friendship, error) {
	if user1 == "" || user2 == "" {
		return nil, errors.New("user input is empty")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	friendship := &database.Friendship{}
	err := c.Pool.QueryRow(
		queryCtx,
		`SELECT
			id, user1, user2, status, created_at, updated_at
		FROM friendships
		WHERE
			( user1 = $1 AND user2 = $2 )
			OR 
			( user1 = $2 AND user2 = $1 )`,
		user1,
		user2,
	).Scan(
		&friendship.Id,
		&friendship.User1,
		&friendship.User2,
		&friendship.Status,
		&friendship.CreatedAt,
		&friendship.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, database.Err_NotFound
		}
		return nil, fmt.Errorf("scanning row: %s", err.Error())
	}
	return friendship, nil
}

func (c *Client) GetFriendshipById(ctx context.Context, friendshipId int32) (*database.Friendship, error) {
	if friendshipId == 0 {
		return nil, errors.New("friendshipId input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	friendship := &database.Friendship{}
	err := c.Pool.QueryRow(
		queryCtx,
		`SELECT
				id, user1, user2, status, created_at, updated_at
			FROM friendships
			WHERE id = $1`,
		friendshipId,
	).Scan(
		&friendship.Id,
		&friendship.User1,
		&friendship.User2,
		&friendship.Status,
		&friendship.CreatedAt,
		&friendship.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, database.Err_NotFound
		}
		return nil, fmt.Errorf("scanning row: %s", err.Error())
	}

	return friendship, nil
}
