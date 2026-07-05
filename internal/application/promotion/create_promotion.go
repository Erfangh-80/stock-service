package promotion

import domainpromotion "stock-service/internal/domain/promotion"

type CreatePromotionUseCase struct {
	repo domainpromotion.Repository
}

func NewCreatePromotionUseCase(repo domainpromotion.Repository) *CreatePromotionUseCase {
	return &CreatePromotionUseCase{repo: repo}
}

func (uc *CreatePromotionUseCase) Execute(title string) (*domainpromotion.Promotion, error) {
	p, err := domainpromotion.NewPromotion(title)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}
	return p, nil
}
