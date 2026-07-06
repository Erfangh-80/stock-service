package productimage

import (
	productimagedomain "stock-service/internal/domain/product_image"
)

type CreateImageInput struct {
	ProductID int32
	FileID    int64
	SortOrder int
}

type CreateImageUseCase struct {
	repo productimagedomain.Repository
}

func NewCreateImageUseCase(repo productimagedomain.Repository) *CreateImageUseCase {
	return &CreateImageUseCase{repo: repo}
}

func (uc *CreateImageUseCase) Execute(input CreateImageInput) (*productimagedomain.ProductImage, error) {
	img, err := productimagedomain.NewProductImage(input.ProductID, input.FileID, input.SortOrder)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(img); err != nil {
		return nil, err
	}
	return img, nil
}
