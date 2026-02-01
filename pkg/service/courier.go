package service

import (
	"context"
	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
)

type CourierService interface {
	ListCouriers(ctx context.Context) ([]models.Courier, error)
}

type courierService struct {
	repo repositoryModels.CourierRepo
}

func NewCourierService(repo repositoryModels.CourierRepo) CourierService {
	return &courierService{repo: repo}
}

func (s *courierService) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	return s.repo.ListCouriers(ctx)
}
