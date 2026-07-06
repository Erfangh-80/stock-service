package brand

import (
	"stock-service/internal/domain/brand"
)

type CreateBrandInput struct {
	Name string
	Slug string
}

type CreateBrandUseCase struct {
	repo brand.Repository
}

func NewCreateBrandUseCase(repo brand.Repository) *CreateBrandUseCase {
	return &CreateBrandUseCase{repo: repo}
}

func (uc *CreateBrandUseCase) Execute(input CreateBrandInput) (*brand.Brand, error) {
	b, err := brand.NewBrand(input.Name, input.Slug)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(b); err != nil {
		return nil, err
	}
	return b, nil
}
