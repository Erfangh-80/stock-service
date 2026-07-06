package category

import (
	"stock-service/internal/domain/category"
)

type UpdateCategoryInput struct {
	ID          int64
	Name        *string
	Slug        *string
	Description *string
	ParentID    *int64
	SortOrder   *int
}

type UpdateCategoryUseCase struct {
	repo category.Repository
}

func NewUpdateCategoryUseCase(repo category.Repository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{repo: repo}
}

func (uc *UpdateCategoryUseCase) Execute(input UpdateCategoryInput) (*category.Category, error) {
	c, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, category.ErrCategoryNotFound
	}

	if input.Name != nil {
		if err := c.UpdateName(*input.Name); err != nil {
			return nil, err
		}
	}
	if input.Slug != nil {
		if err := c.UpdateSlug(*input.Slug); err != nil {
			return nil, err
		}
	}
	if input.Description != nil {
		c.UpdateDescription(*input.Description)
	}
	if input.ParentID != nil {
		c.Reparent(*input.ParentID)
	}
	if input.SortOrder != nil {
		c.UpdateSortOrder(*input.SortOrder)
	}

	if err := uc.repo.Save(c); err != nil {
		return nil, err
	}

	return c, nil
}
