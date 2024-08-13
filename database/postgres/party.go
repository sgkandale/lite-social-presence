package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"socialite/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (c *Client) PutParty(ctx context.Context, party *database.Party) error {
	if party == nil {
		return errors.New("party input is nil")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	// TODO insert membership in transcation
	// TODO check rollback approach
	// start transaction
	tx, err := c.Pool.Begin(queryCtx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %s", err.Error())
	}
	// defer transaction commit or rollback
	defer tx.Rollback(queryCtx)

	log.Printf("new party : %+v", party)

	// insert party
	_, err = tx.Exec(
		queryCtx,
		`INSERT INTO party
			(name, creator, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)`,
		party.Name,
		party.Creator,
		party.CreatedAt,
		party.UpdatedAt,
	)
	if err != nil {
		// duplicate entry check
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return database.Err_DuplicatePrimaryKey
		}
		return fmt.Errorf("inserting party: %s", err.Error())
	}

	// insert party membership
	_, err = tx.Exec(
		queryCtx,
		`INSERT INTO party_members
			(party_name, user_name, status, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5)`,
		party.Name,
		party.Creator,
		database.PartyMembership_Status_Active,
		party.CreatedAt,
		party.UpdatedAt,
	)
	if err != nil {
		// duplicate entry check
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return database.Err_DuplicatePrimaryKey
		}
		return fmt.Errorf("inserting party membership: %s", err.Error())
	}

	// commit transaction
	commitErr := tx.Commit(queryCtx)
	if commitErr != nil {
		return fmt.Errorf("committing transaction: %s", commitErr.Error())
	}

	return nil
}

func (c *Client) GetParty(ctx context.Context, partyName string) (*database.Party, error) {
	if partyName == "" {
		return nil, errors.New("party name is empty")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	party := database.Party{Name: partyName}
	err := c.Pool.QueryRow(
		queryCtx,
		`SELECT
			creator, created_at, updated_at
		FROM
			party
		WHERE
			name = $1`,
		partyName,
	).Scan(
		&party.Creator,
		&party.CreatedAt,
		&party.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, database.Err_NotFound
		}
		return nil, fmt.Errorf("scanning row: %s", err.Error())
	}
	return &party, nil
}

func (c *Client) GetCreatedParties(ctx context.Context, userName string) ([]*database.Party, error) {
	if userName == "" {
		return nil, errors.New("user name is empty")
	}

	queryCtx, cancelQueryCtx := context.WithTimeout(ctx, c.timeout)
	defer cancelQueryCtx()

	rows, err := c.Pool.Query(
		queryCtx,
		`SELECT
			name, creator, created_at, updated_at
		FROM
			party
		WHERE
			creator = $1`,
		userName,
	)
	if err != nil {
		return nil, fmt.Errorf("querying rows: %s", err.Error())
	}
	defer rows.Close()

	createdParties := make([]*database.Party, 0)
	for rows.Next() {
		party := database.Party{}
		err := rows.Scan(
			&party.Name,
			&party.Creator,
			&party.CreatedAt,
			&party.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %s", err.Error())
		}
		createdParties = append(createdParties, &party)
	}
	return createdParties, nil
}
