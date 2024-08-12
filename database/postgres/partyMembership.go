package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"socialite/database"

	"github.com/jackc/pgx/v5/pgconn"
)

func (c *Client) PutPartyMembership(ctx context.Context, membership *database.PartyMembership) error {
	if membership == nil {
		return errors.New("membership input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	// start transaction
	tx, err := c.Pool.Begin(queryCtx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %s", err.Error())
	}
	// defer transaction commit/rollback
	defer func() {
		if err != nil {
			tx.Rollback(queryCtx)
		} else {
			err = tx.Commit(queryCtx)
			if err != nil {
				log.Printf("[ERROR] running add party member transaction: %s", err)
			}
		}
	}()

	// insert party membership
	_, err = tx.Exec(
		queryCtx,
		`INSERT INTO party_members
			(party_name, user_name, status, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5)`,
		membership.PartyName,
		membership.UserName,
		membership.Status,
		membership.CreatedAt,
		membership.UpdatedAt,
	)
	if err != nil {
		// duplicate entry check
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return database.Err_DuplicatePrimaryKey
		}
		return fmt.Errorf("inserting party membership: %s", err.Error())
	}

	// update party table
	_, err = tx.Exec(
		queryCtx,
		`UPDATE party
		SET
			updated_at = $1
		WHERE
			name = $2`,
		time.Now(),
		membership.PartyName,
	)
	if err != nil {
		return fmt.Errorf("updating party updated_at: %s", err.Error())
	}

	return nil
}

func (c *Client) UpdatePartyMembership(ctx context.Context, membership *database.PartyMembership) error {
	if membership == nil {
		return errors.New("membership input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	// start transaction
	tx, err := c.Pool.Begin(queryCtx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %s", err.Error())
	}
	// defer transaction commit/rollback
	defer func() {
		if err != nil {
			tx.Rollback(queryCtx)
		} else {
			err = tx.Commit(queryCtx)
			if err != nil {
				log.Printf("[ERROR] running add party member transaction: %s", err)
			}
		}
	}()

	// update party membership
	_, err = tx.Exec(
		queryCtx,
		`UPDATE party_members
		SET
			status = $1,
			updated_at = $2
		WHERE
			party_name = $3
			AND 
			user_name = $4`,
		membership.Status,
		time.Now(),
		membership.PartyName,
		membership.UserName,
	)
	if err != nil {
		return fmt.Errorf("updating party membership: %s", err.Error())
	}

	// update party table
	_, err = tx.Exec(
		queryCtx,
		`UPDATE party
		SET
			updated_at = $1
		WHERE
			name = $2`,
		time.Now(),
		membership.PartyName,
	)
	if err != nil {
		return fmt.Errorf("updating party updated_at: %s", err.Error())
	}

	return nil
}

func (c *Client) DeletePartyMembership(ctx context.Context, membership *database.PartyMembership) error {
	if membership == nil {
		return errors.New("membership input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	// start transaction
	tx, err := c.Pool.Begin(queryCtx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %s", err.Error())
	}
	// defer transaction commit/rollback
	defer func() {
		if err != nil {
			tx.Rollback(queryCtx)
		} else {
			err = tx.Commit(queryCtx)
			if err != nil {
				log.Printf("[ERROR] running add party member transaction: %s", err)
			}
		}
	}()

	// delete party membership
	_, err = tx.Exec(
		queryCtx,
		`DELETE FROM party_members
		WHERE
			party_name = $1
			AND 
			user_name = $2`,
		membership.PartyName,
		membership.UserName,
	)
	if err != nil {
		return fmt.Errorf("updating party membership: %s", err.Error())
	}

	// update party table
	_, err = tx.Exec(
		queryCtx,
		`UPDATE party
		SET
			updated_at = $1
		WHERE
			name = $2`,
		time.Now(),
		membership.PartyName,
	)
	if err != nil {
		return fmt.Errorf("updating party updated_at: %s", err.Error())
	}

	return nil
}
