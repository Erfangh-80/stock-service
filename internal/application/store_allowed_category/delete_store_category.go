package storeallowedcategory

import (
	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
)

type DeleteStoreCategoryInput struct {
	ID int64
}

type DeleteStoreCategoryUseCase struct {
	repo domainstoreallowedcategory.Repository
}

func NewDeleteStoreCategoryUseCase(repo domainstoreallowedcategory.Repository) *DeleteStoreCategoryUseCase {
	return &DeleteStoreCategoryUseCase{repo: repo}
}

func (uc *DeleteStoreCategoryUseCase) Execute(input DeleteStoreCategoryInput) error {
	sac, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if sac == nil {
		return domainstoreallowedcategory.ErrStoreCategoryNotFound
	}
	return uc.repo.Delete(input.ID)
}
