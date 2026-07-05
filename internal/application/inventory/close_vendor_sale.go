package inventory

import (
	"stock-service/internal/domain/inventory"
)

type CloseVendorSaleInput struct {
	SaleID int64
}

type CloseVendorSaleUseCase interface {
	Execute(input CloseVendorSaleInput) (*inventory.Inventory, error)
}

type closeVendorSaleInteractor struct {
	repo inventory.Repository
}

func NewCloseVendorSaleUseCase(repo inventory.Repository) CloseVendorSaleUseCase {
	return &closeVendorSaleInteractor{repo: repo}
}

func (uc *closeVendorSaleInteractor) Execute(input CloseVendorSaleInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	if err := inv.CloseVendorSale(); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
