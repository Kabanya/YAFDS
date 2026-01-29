package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kabanya/YAFDS/pkg/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCourierPostgresRepository_ListCouriers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewCourierPostgresRepository(db)

	t.Run("success", func(t *testing.T) {
		c1ID := uuid.New()
		c2ID := uuid.New()

		rows := sqlmock.NewRows([]string{"emp_id", "name", "transport_type", "is_active"}).
			AddRow(c1ID, "Courier 1", "bike", true).
			AddRow(c2ID, "Courier 2", "car", false)

		mock.ExpectQuery("SELECT emp_id, name, transport_type, is_active FROM COURIERS").
			WillReturnRows(rows)

		couriers, err := repo.ListCouriers(context.Background())

		assert.NoError(t, err)
		assert.Len(t, couriers, 2)
		assert.Equal(t, c1ID, couriers[0].ID)
		assert.Equal(t, "Courier 1", couriers[0].Name)
		assert.Equal(t, "bike", couriers[0].TransportType)
		assert.True(t, couriers[0].IsActive)
		assert.Equal(t, c2ID, couriers[1].ID)
	})

	t.Run("empty list", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"emp_id", "name", "transport_type", "is_active"})

		mock.ExpectQuery("SELECT emp_id, name, transport_type, is_active FROM COURIERS").
			WillReturnRows(rows)

		couriers, err := repo.ListCouriers(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, couriers)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT emp_id, name, transport_type, is_active FROM COURIERS").
			WillReturnError(errors.New("db error"))

		couriers, err := repo.ListCouriers(context.Background())

		assert.Error(t, err)
		assert.Nil(t, couriers)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"emp_id", "name", "transport_type", "is_active"}).
			AddRow("invalid-uuid", "Courier 1", "bike", true)

		mock.ExpectQuery("SELECT emp_id, name, transport_type, is_active FROM COURIERS").
			WillReturnRows(rows)

		couriers, err := repo.ListCouriers(context.Background())

		assert.Error(t, err)
		assert.Nil(t, couriers)
	})
}
