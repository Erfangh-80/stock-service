package category

import (
	"stock-service/internal/domain/category"
)

type ListCategoriesOutput struct {
	Categories []*category.Category
}

type ListCategoriesUseCase struct {
	repo category.Repository
}

func NewListCategoriesUseCase(repo category.Repository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{repo: repo}
}

func (uc *ListCategoriesUseCase) Execute() (*ListCategoriesOutput, error) {
	categories, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return &ListCategoriesOutput{Categories: categories}, nil
}
