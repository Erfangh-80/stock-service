package storewarehouselink

import (
	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
)

type CreateLinkUseCase struct {
	repo domainstorewarehouselink.Repository
}

func NewCreateLinkUseCase(repo domainstorewarehouselink.Repository) *CreateLinkUseCase {
	return &CreateLinkUseCase{repo: repo}
}

func (uc *CreateLinkUseCase) Execute(storeID, warehouseID int64) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl := domainstorewarehouselink.NewStoreWarehouseLink(storeID, warehouseID)
	if err := uc.repo.Save(swl); err != nil {
		return nil, err
	}
	return swl, nil
}
