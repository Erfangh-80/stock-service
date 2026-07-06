package storewarehouselink

import (
	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
)

type DeleteLinkInput struct {
	ID int64
}

type DeleteLinkUseCase struct {
	repo domainstorewarehouselink.Repository
}

func NewDeleteLinkUseCase(repo domainstorewarehouselink.Repository) *DeleteLinkUseCase {
	return &DeleteLinkUseCase{repo: repo}
}

func (uc *DeleteLinkUseCase) Execute(input DeleteLinkInput) error {
	swl, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if swl == nil {
		return domainstorewarehouselink.ErrLinkNotFound
	}
	return uc.repo.Delete(input.ID)
}
