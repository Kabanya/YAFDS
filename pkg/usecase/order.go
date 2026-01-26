package usecase

import (
	"context"
	"errors"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"

	"github.com/google/uuid"
)

// НА УРОВНЕ USECASE ходим за продуктовыми логикой
var (
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrWalletUnavailable       = errors.New("wallet service unavailable")
	ErrInsufficientFunds       = errors.New("insufficient funds")
)

type WalletClient interface {
	CheckAndDebit(ctx context.Context, walletAddress string, amount float64) (bool, error)
}

type OrderUseCase interface {
	Pay(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID) (models.OrderStatus, error)
	ChangeStatus(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus) (models.OrderStatus, error)
}

type orderUseCase struct {
	repo   repositoryModels.Order
	wallet WalletClient
}
