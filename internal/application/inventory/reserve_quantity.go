package inventory

import (
	"stock-service/internal/domain/inventory"
)

type ReserveQuantityInput struct {
	SaleID    int64
	Quantity  int
}

type ReserveQuantityUseCase interface {
	Execute(input ReserveQuantityInput) (*inventory.Inventory, error)
}

type reserveQuantityInteractor struct {
	repo inventory.Repository
}

func NewReserveQuantityUseCase(repo inventory.Repository) ReserveQuantityUseCase {
	return &reserveQuantityInteractor{repo: repo}
}

func (uc *reserveQuantityInteractor) Execute(input ReserveQuantityInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	if err := inv.ReserveQuantity(input.Quantity); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
