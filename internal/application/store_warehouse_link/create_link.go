package storewarehouselink

import (
	"stock-service/internal/domain/store_warehouse_link"
)

type CreateLinkUseCase struct {
	repo storewarehouselink.Repository
}

func NewCreateLinkUseCase(repo storewarehouselink.Repository) *CreateLinkUseCase {
	return &CreateLinkUseCase{repo: repo}
}

func (uc *CreateLinkUseCase) Execute(storeID, warehouseID int64) (*storewarehouselink.StoreWarehouseLink, error) {
	swl := storewarehouselink.NewStoreWarehouseLink(storeID, warehouseID)
	if err := uc.repo.Save(swl); err != nil {
		return nil, err
	}
	return swl, nil
}
