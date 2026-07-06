package storeallowedcategory

import (
	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
)

type GetStoreCategoryInput struct {
	ID int64
}

type GetStoreCategoryUseCase struct {
	repo domainstoreallowedcategory.Repository
}

func NewGetStoreCategoryUseCase(repo domainstoreallowedcategory.Repository) *GetStoreCategoryUseCase {
	return &GetStoreCategoryUseCase{repo: repo}
}

func (uc *GetStoreCategoryUseCase) Execute(input GetStoreCategoryInput) (*domainstoreallowedcategory.StoreAllowedCategory, error) {
	sac, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if sac == nil {
		return nil, domainstoreallowedcategory.ErrStoreCategoryNotFound
	}
	return sac, nil
}
