package models

import (
	"context"
	pkgModels "github.com/Kabanya/YAFDS/pkg/models"
)

type CourierRepo interface {
	ListCouriers(ctx context.Context) ([]pkgModels.Courier, error)
}
