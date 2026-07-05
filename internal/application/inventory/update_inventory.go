package inventory

import (
	"stock-service/internal/domain/inventory"
)

type UpdateInventoryInput struct {
	SaleID       int64
	InstantQty   int
	ScheduledQty map[string]int
	MinOrderQty  int
	MaxOrderQty  *int
}

type UpdateInventoryUseCase interface {
	Execute(input UpdateInventoryInput) (*inventory.Inventory, error)
}

type updateInventoryInteractor struct {
	repo inventory.Repository
}

func NewUpdateInventoryUseCase(repo inventory.Repository) UpdateInventoryUseCase {
	return &updateInventoryInteractor{repo: repo}
}

func (uc *updateInventoryInteractor) Execute(input UpdateInventoryInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if err := inv.UpdateInventory(input.InstantQty, input.ScheduledQty, input.MinOrderQty, input.MaxOrderQty); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
