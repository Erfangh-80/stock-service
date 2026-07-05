package inventory

import (
	"stock-service/internal/domain/inventory"
)

type CheckLowStockInput struct {
	SaleID    int64
	Threshold int
}

type CheckLowStockOutput struct {
	IsLow     bool
	CurrentQty int
}

type CheckLowStockUseCase interface {
	Execute(input CheckLowStockInput) (*CheckLowStockOutput, error)
}

type checkLowStockInteractor struct {
	repo inventory.Repository
}

func NewCheckLowStockUseCase(repo inventory.Repository) CheckLowStockUseCase {
	return &checkLowStockInteractor{repo: repo}
}

func (uc *checkLowStockInteractor) Execute(input CheckLowStockInput) (*CheckLowStockOutput, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	return &CheckLowStockOutput{
		IsLow:      inv.HasLowStock(input.Threshold),
		CurrentQty: inv.InstantQty,
	}, nil
}
