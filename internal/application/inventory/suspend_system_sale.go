package inventory

import (
	"stock-service/internal/domain/inventory"
)

type SuspendSystemSaleInput struct {
	SaleID int64
}

type SuspendSystemSaleUseCase interface {
	Execute(input SuspendSystemSaleInput) (*inventory.Inventory, error)
}

type suspendSystemSaleInteractor struct {
	repo inventory.Repository
}

func NewSuspendSystemSaleUseCase(repo inventory.Repository) SuspendSystemSaleUseCase {
	return &suspendSystemSaleInteractor{repo: repo}
}

func (uc *suspendSystemSaleInteractor) Execute(input SuspendSystemSaleInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	if err := inv.SuspendSystemSale(); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
