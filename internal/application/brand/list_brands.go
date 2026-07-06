package brand

import (
	"stock-service/internal/domain/brand"
)

type ListBrandsOutput struct {
	Brands []*brand.Brand
}

type ListBrandsUseCase struct {
	repo brand.Repository
}

func NewListBrandsUseCase(repo brand.Repository) *ListBrandsUseCase {
	return &ListBrandsUseCase{repo: repo}
}

func (uc *ListBrandsUseCase) Execute() (*ListBrandsOutput, error) {
	brands, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return &ListBrandsOutput{Brands: brands}, nil
}
