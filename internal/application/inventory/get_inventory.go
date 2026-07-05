package inventory

import (
	"stock-service/internal/domain/inventory"
)

type GetInventoryInput struct {
	ID int64
}

type GetInventoryUseCase interface {
	Execute(input GetInventoryInput) (*inventory.Inventory, error)
}

type getInventoryInteractor struct {
	repo inventory.Repository
}

func NewGetInventoryUseCase(repo inventory.Repository) GetInventoryUseCase {
	return &getInventoryInteractor{repo: repo}
}

func (uc *getInventoryInteractor) Execute(input GetInventoryInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}
	return inv, nil
}
