package promotion

import (
	domainpromotion "stock-service/internal/domain/promotion"
)

type DeletePromotionInput struct {
	ID int64
}

type DeletePromotionUseCase struct {
	repo domainpromotion.Repository
}

func NewDeletePromotionUseCase(repo domainpromotion.Repository) *DeletePromotionUseCase {
	return &DeletePromotionUseCase{repo: repo}
}

func (uc *DeletePromotionUseCase) Execute(input DeletePromotionInput) error {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if p == nil {
		return domainpromotion.ErrPromotionNotFound
	}
	return uc.repo.Delete(input.ID)
}
