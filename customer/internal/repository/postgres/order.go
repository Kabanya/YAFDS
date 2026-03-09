package repository

import (
	"context"
	"database/sql"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"

	"github.com/google/uuid"
)

type OrderUserRepo interface {
	CreateOrder(ctx context.Context, order repositoryModels.Filter) (models.Order, error)                                        // Customer
	CreateOrderWithItems(ctx context.Context, order models.Order, items []repositoryModels.OrderItemInput) (models.Order, error) // Customer
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)                                          // Customer
	PayOrder(ctx context.Context, orderID uuid.UUID) error                                                                       // Customer
}

type postgresUserRepo struct {
	customerDB *sql.DB
	orderDB    *sql.DB
}

func NewPostgresUserRepo(customerDB, orderDB *sql.DB) OrderUserRepo {
	return &postgresUserRepo{customerDB: customerDB, orderDB: orderDB}
}

// CreateOrder implements [OrderUserRepo].
func (p *postgresUserRepo) CreateOrder(ctx context.Context, order repositoryModels.Filter) (models.Order, error) {
	panic("unimplemented")
}

// CreateOrderWithItems implements [OrderUserRepo].
func (p *postgresUserRepo) CreateOrderWithItems(ctx context.Context, order models.Order, items []repositoryModels.OrderItemInput) (models.Order, error) {
	panic("unimplemented")
}

// GetCustomerWalletAddress implements [OrderUserRepo].
func (p *postgresUserRepo) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	panic("unimplemented")
}

// PayOrder implements [OrderUserRepo].
func (p *postgresUserRepo) PayOrder(ctx context.Context, orderID uuid.UUID) error {
	panic("unimplemented")
}
