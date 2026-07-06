package category

import (
	"stock-service/internal/domain/category"
)

type CreateCategoryInput struct {
	Name        string
	Slug        string
	ParentID    *int64
	Description *string
}

type CreateCategoryUseCase struct {
	repo category.Repository
}

func NewCreateCategoryUseCase(repo category.Repository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{repo: repo}
}

func (uc *CreateCategoryUseCase) Execute(input CreateCategoryInput) (*category.Category, error) {
	c, err := category.NewCategory(input.Name, input.Slug, input.ParentID)
	if err != nil {
		return nil, err
	}

	if input.Description != nil {
		c.UpdateDescription(*input.Description)
	}

	if err := uc.repo.Save(c); err != nil {
		return nil, err
	}

	return c, nil
}
