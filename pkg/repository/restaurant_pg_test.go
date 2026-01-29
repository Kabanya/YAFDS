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

func TestRestaurantPostgresRepository_ListRestaurants(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewRestaurantPostgresRepository(db)

	t.Run("success", func(t *testing.T) {
		res1ID := uuid.New()
		res2ID := uuid.New()

		rows := sqlmock.NewRows([]string{"emp_id", "name", "address", "status"}).
			AddRow(res1ID, "Restaurant 1", "Address 1", true).
			AddRow(res2ID, "Restaurant 2", "Address 2", false)

		mock.ExpectQuery("SELECT emp_id, name, address, status FROM RESTAURANTS").
			WillReturnRows(rows)

		restaurants, err := repo.ListRestaurants(context.Background())

		assert.NoError(t, err)
		assert.Len(t, restaurants, 2)
		assert.Equal(t, res1ID, restaurants[0].ID)
		assert.Equal(t, "Restaurant 1", restaurants[0].Name)
		assert.Equal(t, "Address 1", restaurants[0].Address)
		assert.True(t, restaurants[0].Status)
		assert.Equal(t, res2ID, restaurants[1].ID)
	})

	t.Run("empty list", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"emp_id", "name", "address", "status"})

		mock.ExpectQuery("SELECT emp_id, name, address, status FROM RESTAURANTS").
			WillReturnRows(rows)

		restaurants, err := repo.ListRestaurants(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, restaurants)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT emp_id, name, address, status FROM RESTAURANTS").
			WillReturnError(errors.New("db error"))

		restaurants, err := repo.ListRestaurants(context.Background())

		assert.Error(t, err)
		assert.Nil(t, restaurants)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"emp_id", "name", "address", "status"}).
			AddRow("invalid-uuid", "Restaurant 1", "Address 1", true)

		mock.ExpectQuery("SELECT emp_id, name, address, status FROM RESTAURANTS").
			WillReturnRows(rows)

		restaurants, err := repo.ListRestaurants(context.Background())

		assert.Error(t, err)
		assert.Nil(t, restaurants)
	})
}
