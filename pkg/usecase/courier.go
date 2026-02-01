package usecase

import (
	"context"
	"github.com/Kabanya/YAFDS/pkg/models"
	"github.com/Kabanya/YAFDS/pkg/service"
)

type CourierUseCase interface {
	ListCouriers(ctx context.Context) ([]models.Courier, error)
}

type courierUseCase struct {
	svc service.CourierService
}

func NewCourierUseCase(svc service.CourierService) CourierUseCase {
	return &courierUseCase{svc: svc}
}

func (u *courierUseCase) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	return u.svc.ListCouriers(ctx)
}
