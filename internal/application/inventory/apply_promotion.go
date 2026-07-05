package inventory

import (
	"time"

	"stock-service/internal/domain/inventory"
)

type ApplyPromotionInput struct {
	SaleID      int64
	PromotionID int64
	FinalPrice  float64
	StartAt     time.Time
	EndAt       time.Time
}

type ApplyPromotionUseCase interface {
	Execute(input ApplyPromotionInput) (*inventory.Inventory, error)
}

type applyPromotionInteractor struct {
	repo inventory.Repository
}

func NewApplyPromotionUseCase(repo inventory.Repository) ApplyPromotionUseCase {
	return &applyPromotionInteractor{repo: repo}
}

func (uc *applyPromotionInteractor) Execute(input ApplyPromotionInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if err := inv.ApplyPromotion(input.PromotionID, input.FinalPrice, input.StartAt, input.EndAt); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
