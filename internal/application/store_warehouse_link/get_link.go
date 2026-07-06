package storewarehouselink

import (
	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
)

type GetLinkInput struct {
	ID int64
}

type GetLinkUseCase struct {
	repo domainstorewarehouselink.Repository
}

func NewGetLinkUseCase(repo domainstorewarehouselink.Repository) *GetLinkUseCase {
	return &GetLinkUseCase{repo: repo}
}

func (uc *GetLinkUseCase) Execute(input GetLinkInput) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if swl == nil {
		return nil, domainstorewarehouselink.ErrLinkNotFound
	}
	return swl, nil
}
