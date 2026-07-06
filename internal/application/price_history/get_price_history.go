package pricehistory

import (
	pricehistorydomain "stock-service/internal/domain/price_history"
)

type GetPriceHistoryInput struct {
	ProductID int32
}

type GetPriceHistoryOutput struct {
	History []*pricehistorydomain.PriceHistory
}

type GetPriceHistoryUseCase struct {
	repo pricehistorydomain.Repository
}

func NewGetPriceHistoryUseCase(repo pricehistorydomain.Repository) *GetPriceHistoryUseCase {
	return &GetPriceHistoryUseCase{repo: repo}
}

func (uc *GetPriceHistoryUseCase) Execute(input GetPriceHistoryInput) (*GetPriceHistoryOutput, error) {
	history, err := uc.repo.FindByProductID(input.ProductID)
	if err != nil {
		return nil, err
	}
	return &GetPriceHistoryOutput{History: history}, nil
}
