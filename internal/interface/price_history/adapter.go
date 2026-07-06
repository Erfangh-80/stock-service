package pricehistoryinterface

import (
	apppricehistory "stock-service/internal/application/price_history"
	"stock-service/internal/domain/price_history"
	iface "stock-service/internal/interface"
)

type PriceHistoryResponse struct {
	ID          int64   `json:"id"`
	ProductID   int32   `json:"product_id"`
	OldPrice    float64 `json:"old_price"`
	NewPrice    float64 `json:"new_price"`
	ChangedBy   string  `json:"changed_by"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type CreatePriceHistoryParams struct {
	ProductID   int32
	OldPrice    float64
	NewPrice    float64
	ChangedBy   string
	Description *string
}

type Adapter struct {
	create apppricehistory.CreatePriceHistoryUseCase
	list   apppricehistory.GetPriceHistoryUseCase
}

func NewAdapter(
	create apppricehistory.CreatePriceHistoryUseCase,
	list apppricehistory.GetPriceHistoryUseCase,
) *Adapter {
	return &Adapter{create: create, list: list}
}

func (a *Adapter) Create(params CreatePriceHistoryParams) (*PriceHistoryResponse, error) {
	result, err := a.create.Execute(apppricehistory.CreatePriceHistoryInput{
		ProductID: params.ProductID, OldPrice: params.OldPrice,
		NewPrice: params.NewPrice, ChangedBy: params.ChangedBy,
		Description: params.Description,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) List(productID int32) ([]*PriceHistoryResponse, error) {
	result, err := a.list.Execute(apppricehistory.GetPriceHistoryInput{ProductID: productID})
	if err != nil {
		return nil, mapError(err)
	}
	resp := make([]*PriceHistoryResponse, len(result.History))
	for i, ph := range result.History {
		resp[i] = toResponse(ph)
	}
	return resp, nil
}

func mapError(err error) error {
	switch err {
	case pricehistory.ErrInvalidProductID, pricehistory.ErrInvalidPrice, pricehistory.ErrChangedByRequired:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(ph *pricehistory.PriceHistory) *PriceHistoryResponse {
	return &PriceHistoryResponse{
		ID: ph.ID, ProductID: ph.ProductID, OldPrice: ph.OldPrice,
		NewPrice: ph.NewPrice, ChangedBy: ph.ChangedBy,
		Description: ph.Description,
		CreatedAt:   ph.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
