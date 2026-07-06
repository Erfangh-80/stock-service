package pricehistory

import (
	pricehistorydomain "stock-service/internal/domain/price_history"
)

type CreatePriceHistoryInput struct {
	ProductID   int32
	OldPrice    float64
	NewPrice    float64
	ChangedBy   string
	Description *string
}

type CreatePriceHistoryUseCase struct {
	repo pricehistorydomain.Repository
}

func NewCreatePriceHistoryUseCase(repo pricehistorydomain.Repository) *CreatePriceHistoryUseCase {
	return &CreatePriceHistoryUseCase{repo: repo}
}

func (uc *CreatePriceHistoryUseCase) Execute(input CreatePriceHistoryInput) (*pricehistorydomain.PriceHistory, error) {
	ph, err := pricehistorydomain.NewPriceHistory(input.ProductID, input.OldPrice, input.NewPrice, input.ChangedBy, input.Description)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ph); err != nil {
		return nil, err
	}
	return ph, nil
}
