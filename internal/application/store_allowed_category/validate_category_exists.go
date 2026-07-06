package storeallowedcategory

import (
	domaincategory "stock-service/internal/domain/category"
	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
)

type ValidateCategoryExistsInput struct {
	CategoryID int64
}

type ValidateCategoryExistsUseCase struct {
	storeCategoryRepo domainstoreallowedcategory.Repository
	categoryRepo      domaincategory.Repository
}

func NewValidateCategoryExistsUseCase(
	storeCategoryRepo domainstoreallowedcategory.Repository,
	categoryRepo domaincategory.Repository,
) *ValidateCategoryExistsUseCase {
	return &ValidateCategoryExistsUseCase{
		storeCategoryRepo: storeCategoryRepo,
		categoryRepo:      categoryRepo,
	}
}

func (uc *ValidateCategoryExistsUseCase) Execute(input ValidateCategoryExistsInput) error {
	cat, err := uc.categoryRepo.FindByID(input.CategoryID)
	if err != nil {
		return err
	}
	if cat == nil {
		return domainstoreallowedcategory.ErrCategoryNotFound
	}
	return nil
}
