package storeallowedcategory

import (
	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
)

type ListStoreCategoriesInput struct {
	StoreID *int64
	Page    int
	Limit   int
}

type ListStoreCategoriesOutput struct {
	Categories []*domainstoreallowedcategory.StoreAllowedCategory
	Total      int
	Page       int
	Limit      int
}

type ListStoreCategoriesUseCase struct {
	repo domainstoreallowedcategory.Repository
}

func NewListStoreCategoriesUseCase(repo domainstoreallowedcategory.Repository) *ListStoreCategoriesUseCase {
	return &ListStoreCategoriesUseCase{repo: repo}
}

func (uc *ListStoreCategoriesUseCase) Execute(input ListStoreCategoriesInput) (*ListStoreCategoriesOutput, error) {
	filter := domainstoreallowedcategory.StoreCategoryFilter{
		StoreID: input.StoreID,
		Page:    input.Page,
		Limit:   input.Limit,
	}
	categories, total, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return &ListStoreCategoriesOutput{
		Categories: categories,
		Total:      total,
		Page:       input.Page,
		Limit:      input.Limit,
	}, nil
}
