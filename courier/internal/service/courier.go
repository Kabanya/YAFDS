package service

import (
	"context"

	"github.com/Kabanya/YAFDS/courier/internal/domain"
)

type CourierService interface {
	ListCouriers(ctx context.Context) ([]domain.Courier, error)
}

type courierService struct {
	repo domain.CourierRepo
}

func NewCourierService(repo domain.CourierRepo) CourierService {
	return &courierService{repo: repo}
}

func (s *courierService) ListCouriers(ctx context.Context) ([]domain.Courier, error) {
	return s.repo.ListCouriers(ctx)
}
