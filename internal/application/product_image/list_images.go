package productimage

import (
	productimagedomain "stock-service/internal/domain/product_image"
)

type ListImagesInput struct {
	ProductID int32
}

type ListImagesOutput struct {
	Images []*productimagedomain.ProductImage
}

type ListImagesUseCase struct {
	repo productimagedomain.Repository
}

func NewListImagesUseCase(repo productimagedomain.Repository) *ListImagesUseCase {
	return &ListImagesUseCase{repo: repo}
}

func (uc *ListImagesUseCase) Execute(input ListImagesInput) (*ListImagesOutput, error) {
	images, err := uc.repo.FindByProductID(input.ProductID)
	if err != nil {
		return nil, err
	}
	return &ListImagesOutput{Images: images}, nil
}
