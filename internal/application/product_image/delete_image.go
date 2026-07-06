package productimage

import (
	productimagedomain "stock-service/internal/domain/product_image"
)

type DeleteImageInput struct {
	ID int64
}

type DeleteImageUseCase struct {
	repo productimagedomain.Repository
}

func NewDeleteImageUseCase(repo productimagedomain.Repository) *DeleteImageUseCase {
	return &DeleteImageUseCase{repo: repo}
}

func (uc *DeleteImageUseCase) Execute(input DeleteImageInput) error {
	img, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if img == nil {
		return productimagedomain.ErrImageNotFound
	}
	return uc.repo.Delete(input.ID)
}
