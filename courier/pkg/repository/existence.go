package pkg

import (
	"context"
	"database/sql"
	"errors"

	pkgErrors "github.com/Kabanya/YAFDS/pkg/errors"
	"github.com/google/uuid"
)

func EnsureCourierExists(ctx context.Context, db *sql.DB, courierID uuid.UUID) error {
	const query = "SELECT id FROM COURIERS WHERE id = $1"

	var id uuid.UUID
	if err := db.QueryRowContext(ctx, query, courierID).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkgErrors.ErrCourierNotFound
		}
		return err
	}

	return nil
}
