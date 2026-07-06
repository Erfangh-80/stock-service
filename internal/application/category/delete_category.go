package category

import (
	"stock-service/internal/domain/category"
)

type DeleteCategoryInput struct {
	ID int64
}

type DeleteCategoryUseCase struct {
	repo category.Repository
}

func NewDeleteCategoryUseCase(repo category.Repository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{repo: repo}
}

func (uc *DeleteCategoryUseCase) Execute(input DeleteCategoryInput) error {
	c, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if c == nil {
		return category.ErrCategoryNotFound
	}
	return uc.repo.Delete(input.ID)
}
