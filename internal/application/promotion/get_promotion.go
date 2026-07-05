package promotion

import domainpromotion "stock-service/internal/domain/promotion"

type GetPromotionInput struct {
	ID int64
}

type GetPromotionUseCase struct {
	repo domainpromotion.Repository
}

func NewGetPromotionUseCase(repo domainpromotion.Repository) *GetPromotionUseCase {
	return &GetPromotionUseCase{repo: repo}
}

func (uc *GetPromotionUseCase) Execute(input GetPromotionInput) (*domainpromotion.Promotion, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, domainpromotion.ErrPromotionNotFound
	}
	return p, nil
}
