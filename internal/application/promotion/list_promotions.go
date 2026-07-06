package promotion

import (
	domainpromotion "stock-service/internal/domain/promotion"
)

type ListPromotionsInput struct {
	Status       *domainpromotion.PromotionStatus
	DiscountType *domainpromotion.DiscountType
	Search       *string
	Page         int
	Limit        int
}

type ListPromotionsOutput struct {
	Promotions []*domainpromotion.Promotion
	Total      int
	Page       int
	Limit      int
}

type ListPromotionsUseCase struct {
	repo domainpromotion.Repository
}

func NewListPromotionsUseCase(repo domainpromotion.Repository) *ListPromotionsUseCase {
	return &ListPromotionsUseCase{repo: repo}
}

func (uc *ListPromotionsUseCase) Execute(input ListPromotionsInput) (*ListPromotionsOutput, error) {
	filter := domainpromotion.PromotionFilter{
		Status:       input.Status,
		DiscountType: input.DiscountType,
		Search:       input.Search,
		Page:         input.Page,
		Limit:        input.Limit,
	}

	items, total, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return &ListPromotionsOutput{
		Promotions: items,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
	}, nil
}
