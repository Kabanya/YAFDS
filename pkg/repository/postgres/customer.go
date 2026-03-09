package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func EnsureCustomerExists(ctx context.Context, db *sql.DB, customerID uuid.UUID) error {
	const query = "SELECT id FROM CUSTOMERS WHERE id = $1"
	var id uuid.UUID
	if err := db.QueryRowContext(ctx, query, customerID).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("customer not found")
		}
		return err
	}
	return nil
}
