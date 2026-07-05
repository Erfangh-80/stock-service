package inventory

import (
	"stock-service/internal/domain/inventory"
)

type DeleteInventoryInput struct {
	SaleID int64
}

type DeleteInventoryUseCase interface {
	Execute(input DeleteInventoryInput) error
}

type deleteInventoryInteractor struct {
	repo inventory.Repository
}

func NewDeleteInventoryUseCase(repo inventory.Repository) DeleteInventoryUseCase {
	return &deleteInventoryInteractor{repo: repo}
}

func (uc *deleteInventoryInteractor) Execute(input DeleteInventoryInput) error {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return err
	}
	if inv == nil {
		return inventory.ErrInventoryNotFound
	}
	return uc.repo.Delete(input.SaleID)
}
