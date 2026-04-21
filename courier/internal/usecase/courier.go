package usecase

import (
	"context"

	"github.com/Kabanya/YAFDS/courier/internal/domain"
	"github.com/Kabanya/YAFDS/courier/internal/service"
)

type CourierUseCase interface {
	ListCouriers(ctx context.Context) ([]domain.Courier, error)
}

type courierUseCase struct {
	svc service.CourierService
}

func NewCourierUseCase(svc service.CourierService) CourierUseCase {
	return &courierUseCase{svc: svc}
}

func (u *courierUseCase) ListCouriers(ctx context.Context) ([]domain.Courier, error) {
	return u.svc.ListCouriers(ctx)
}
