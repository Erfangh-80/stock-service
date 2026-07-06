package brand

import (
	"stock-service/internal/domain/brand"
)

type UpdateBrandInput struct {
	ID   int64
	Name *string
	Slug *string
}

type UpdateBrandUseCase struct {
	repo brand.Repository
}

func NewUpdateBrandUseCase(repo brand.Repository) *UpdateBrandUseCase {
	return &UpdateBrandUseCase{repo: repo}
}

func (uc *UpdateBrandUseCase) Execute(input UpdateBrandInput) (*brand.Brand, error) {
	b, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, brand.ErrBrandNotFound
	}

	if input.Name != nil {
		if err := b.UpdateName(*input.Name); err != nil {
			return nil, err
		}
	}
	if input.Slug != nil {
		if err := b.UpdateSlug(*input.Slug); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Save(b); err != nil {
		return nil, err
	}

	return b, nil
}
