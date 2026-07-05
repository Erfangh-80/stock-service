package inventory

import (
	"stock-service/internal/domain/inventory"
)

type ReleaseQuantityInput struct {
	SaleID    int64
	Quantity  int
}

type ReleaseQuantityUseCase interface {
	Execute(input ReleaseQuantityInput) (*inventory.Inventory, error)
}

type releaseQuantityInteractor struct {
	repo inventory.Repository
}

func NewReleaseQuantityUseCase(repo inventory.Repository) ReleaseQuantityUseCase {
	return &releaseQuantityInteractor{repo: repo}
}

func (uc *releaseQuantityInteractor) Execute(input ReleaseQuantityInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	if err := inv.ReleaseQuantity(input.Quantity); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
