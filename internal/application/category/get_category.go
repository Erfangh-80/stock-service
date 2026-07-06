package category

import (
	"stock-service/internal/domain/category"
)

type GetCategoryInput struct {
	ID int64
}

type GetCategoryUseCase struct {
	repo category.Repository
}

func NewGetCategoryUseCase(repo category.Repository) *GetCategoryUseCase {
	return &GetCategoryUseCase{repo: repo}
}

func (uc *GetCategoryUseCase) Execute(input GetCategoryInput) (*category.Category, error) {
	c, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, category.ErrCategoryNotFound
	}
	return c, nil
}
