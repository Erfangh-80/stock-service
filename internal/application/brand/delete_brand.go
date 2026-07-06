package brand

import (
	"stock-service/internal/domain/brand"
)

type DeleteBrandInput struct {
	ID int64
}

type DeleteBrandUseCase struct {
	repo brand.Repository
}

func NewDeleteBrandUseCase(repo brand.Repository) *DeleteBrandUseCase {
	return &DeleteBrandUseCase{repo: repo}
}

func (uc *DeleteBrandUseCase) Execute(input DeleteBrandInput) error {
	b, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if b == nil {
		return brand.ErrBrandNotFound
	}
	return uc.repo.Delete(input.ID)
}
