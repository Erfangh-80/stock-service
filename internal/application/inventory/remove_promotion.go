package inventory

import (
	"stock-service/internal/domain/inventory"
)

type RemovePromotionInput struct {
	SaleID int64
}

type RemovePromotionUseCase interface {
	Execute(input RemovePromotionInput) (*inventory.Inventory, error)
}

type removePromotionInteractor struct {
	repo inventory.Repository
}

func NewRemovePromotionUseCase(repo inventory.Repository) RemovePromotionUseCase {
	return &removePromotionInteractor{repo: repo}
}

func (uc *removePromotionInteractor) Execute(input RemovePromotionInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if err := inv.RemovePromotion(); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
