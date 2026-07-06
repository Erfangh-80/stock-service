package storewarehouselink

import (
	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
)

type ChangeRelationInput struct {
	LinkID       int64
	RelationType domainstorewarehouselink.RelationType
}

type ChangeRelationUseCase struct {
	repo domainstorewarehouselink.Repository
}

func NewChangeRelationUseCase(repo domainstorewarehouselink.Repository) *ChangeRelationUseCase {
	return &ChangeRelationUseCase{repo: repo}
}

func (uc *ChangeRelationUseCase) Execute(input ChangeRelationInput) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl, err := uc.repo.FindByID(input.LinkID)
	if err != nil {
		return nil, err
	}
	if swl == nil {
		return nil, domainstorewarehouselink.ErrLinkNotFound
	}
	if err := swl.ChangeRelationType(input.RelationType); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(swl); err != nil {
		return nil, err
	}
	return swl, nil
}
