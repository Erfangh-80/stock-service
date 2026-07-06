package brand

import (
	"stock-service/internal/domain/brand"
)

type GetBrandInput struct {
	ID int64
}

type GetBrandUseCase struct {
	repo brand.Repository
}

func NewGetBrandUseCase(repo brand.Repository) *GetBrandUseCase {
	return &GetBrandUseCase{repo: repo}
}

func (uc *GetBrandUseCase) Execute(input GetBrandInput) (*brand.Brand, error) {
	b, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, brand.ErrBrandNotFound
	}
	return b, nil
}
