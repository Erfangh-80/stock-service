package inventory

import (
	"stock-service/internal/domain/inventory"
)

type SuspendVendorSaleInput struct {
	SaleID int64
}

type SuspendVendorSaleUseCase interface {
	Execute(input SuspendVendorSaleInput) (*inventory.Inventory, error)
}

type suspendVendorSaleInteractor struct {
	repo inventory.Repository
}

func NewSuspendVendorSaleUseCase(repo inventory.Repository) SuspendVendorSaleUseCase {
	return &suspendVendorSaleInteractor{repo: repo}
}

func (uc *suspendVendorSaleInteractor) Execute(input SuspendVendorSaleInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	if err := inv.SuspendVendorSale(); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
