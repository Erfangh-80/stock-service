package inventory

import (
	"stock-service/internal/domain/inventory"
)

type CloseSystemSaleInput struct {
	SaleID int64
}

type CloseSystemSaleUseCase interface {
	Execute(input CloseSystemSaleInput) (*inventory.Inventory, error)
}

type closeSystemSaleInteractor struct {
	repo inventory.Repository
}

func NewCloseSystemSaleUseCase(repo inventory.Repository) CloseSystemSaleUseCase {
	return &closeSystemSaleInteractor{repo: repo}
}

func (uc *closeSystemSaleInteractor) Execute(input CloseSystemSaleInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	if err := inv.CloseSystemSale(); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
