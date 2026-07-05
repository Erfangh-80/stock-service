package inventory

import (
	"stock-service/internal/domain/inventory"
	productdomain "stock-service/internal/domain/product"
)

type CreateInventoryInput struct {
	StoreID     int64
	WarehouseID int64
	ProductID   int32
	BasePrice   float64
}

type CreateInventoryUseCase interface {
	Execute(input CreateInventoryInput) (*inventory.Inventory, error)
}

type createProductSaleInteractor struct {
	invRepo     inventory.Repository
	productRepo productdomain.Repository
}

func NewCreateInventoryUseCase(invRepo inventory.Repository, productRepo productdomain.Repository) CreateInventoryUseCase {
	return &createProductSaleInteractor{invRepo: invRepo, productRepo: productRepo}
}

func (uc *createProductSaleInteractor) Execute(input CreateInventoryInput) (*inventory.Inventory, error) {
	product, err := uc.productRepo.FindByID(input.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, productdomain.ErrProductNotFound
	}

	inv, err := inventory.NewInventory(input.StoreID, input.WarehouseID, input.ProductID, input.BasePrice)
	if err != nil {
		return nil, err
	}
	if err := uc.invRepo.Save(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
