package usecase

import (
	"context"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type PurchaseRepository interface {
	FindPurchases(ctx context.Context, filter entity.PurchaseFilter) ([]entity.Purchase, error)
	GetMetrics(ctx context.Context, filter entity.PurchaseFilter) (entity.PurchaseMetrics, error)
}

type PurchaseUseCase struct {
	repo PurchaseRepository
}

func NewPurchaseUseCase(repo PurchaseRepository) *PurchaseUseCase {
	return &PurchaseUseCase{repo: repo}
}

func (uc *PurchaseUseCase) GetPurchases(ctx context.Context, filter entity.PurchaseFilter) ([]entity.Purchase, error) {
	return uc.repo.FindPurchases(ctx, filter)
}

func (uc *PurchaseUseCase) GetMetrics(ctx context.Context, filter entity.PurchaseFilter) (entity.PurchaseMetrics, error) {
	return uc.repo.GetMetrics(ctx, filter)
}
