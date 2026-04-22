package pkg

import (
	"context"
	"database/sql"
	"errors"

	pkgErrors "github.com/Kabanya/YAFDS/pkg/errors"
	"github.com/google/uuid"
)

func EnsureCustomerExists(ctx context.Context, db *sql.DB, customerID uuid.UUID) error {
	const query = "SELECT id FROM CUSTOMERS WHERE id = $1"

	var id uuid.UUID
	if err := db.QueryRowContext(ctx, query, customerID).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkgErrors.ErrCustomerNotFound
		}
		return err
	}

	return nil
}
