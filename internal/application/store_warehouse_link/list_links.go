package storewarehouselink

import (
	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
)

type ListLinksInput struct {
	StoreID     *int64
	WarehouseID *int64
	Page        int
	Limit       int
}

type ListLinksOutput struct {
	Links []*domainstorewarehouselink.StoreWarehouseLink
	Total int
	Page  int
	Limit int
}

type ListLinksUseCase struct {
	repo domainstorewarehouselink.Repository
}

func NewListLinksUseCase(repo domainstorewarehouselink.Repository) *ListLinksUseCase {
	return &ListLinksUseCase{repo: repo}
}

func (uc *ListLinksUseCase) Execute(input ListLinksInput) (*ListLinksOutput, error) {
	filter := domainstorewarehouselink.WarehouseLinkFilter{
		StoreID:     input.StoreID,
		WarehouseID: input.WarehouseID,
		Page:        input.Page,
		Limit:       input.Limit,
	}
	links, total, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return &ListLinksOutput{
		Links: links,
		Total: total,
		Page:  input.Page,
		Limit: input.Limit,
	}, nil
}
