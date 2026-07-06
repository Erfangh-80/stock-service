package salescommissioninterface

import (
	"errors"
	"time"

	domain "stock-service/internal/domain/sales_commission"
	iface "stock-service/internal/interface"
)

type SalesCommissionOutput struct {
	ID                       int64
	InventoryID              int64
	CategoryCommissionRuleID int64
	SaleModel                string
	RatePercent              float64
	MinQty                   *int
	MinPrice                 float64
	MaxPrice                 *float64
	CreatedAt                time.Time
}

type CreateSalesCommissionParams struct {
	InventoryID              int64
	CategoryCommissionRuleID int64
	SaleModel                string
	RatePercent              float64
	MinPrice                 float64
}

type UpdateMaxPriceParams struct {
	CommissionID int64
	MaxPrice     float64
}

type UpdateMinQtyParams struct {
	CommissionID int64
	MinQty       int
}

type createSalesCommissionUseCase interface {
	Execute(inventoryID, categoryCommissionRuleID int64, saleModel domain.SaleModel, ratePercent, minPrice float64) (*domain.SalesCommission, error)
}

type updateMaxPriceUseCase interface {
	Execute(commissionID int64, maxPrice float64) error
}

type updateMinQtyUseCase interface {
	Execute(commissionID int64, minQty int) error
}

type Adapter struct {
	create          createSalesCommissionUseCase
	updateMaxPrice  updateMaxPriceUseCase
	updateMinQty    updateMinQtyUseCase
}

func NewAdapter(
	create createSalesCommissionUseCase,
	updateMaxPrice updateMaxPriceUseCase,
	updateMinQty updateMinQtyUseCase,
) *Adapter {
	return &Adapter{
		create:          create,
		updateMaxPrice:  updateMaxPrice,
		updateMinQty:    updateMinQty,
	}
}

func (a *Adapter) Create(params CreateSalesCommissionParams) (*SalesCommissionOutput, error) {
	result, err := a.create.Execute(
		params.InventoryID, params.CategoryCommissionRuleID,
		domain.SaleModel(params.SaleModel), params.RatePercent, params.MinPrice,
	)
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) UpdateMaxPrice(params UpdateMaxPriceParams) error {
	err := a.updateMaxPrice.Execute(params.CommissionID, params.MaxPrice)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) UpdateMinQty(params UpdateMinQtyParams) error {
	err := a.updateMinQty.Execute(params.CommissionID, params.MinQty)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidRatePercent),
		errors.Is(err, domain.ErrInvalidMinPrice),
		errors.Is(err, domain.ErrInvalidMaxPrice),
		errors.Is(err, domain.ErrInvalidMinQty):
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toOutput(sc *domain.SalesCommission) *SalesCommissionOutput {
	var minQty *int
	if sc.MinQty != nil {
		v := *sc.MinQty
		minQty = &v
	}
	var maxPrice *float64
	if sc.MaxPrice != nil {
		v := *sc.MaxPrice
		maxPrice = &v
	}
	return &SalesCommissionOutput{
		ID:                       sc.ID,
		InventoryID:              sc.InventoryID,
		CategoryCommissionRuleID: sc.CategoryCommissionRuleID,
		SaleModel:                string(sc.SaleModel),
		RatePercent:              sc.RatePercent,
		MinQty:                   minQty,
		MinPrice:                 sc.MinPrice,
		MaxPrice:                 maxPrice,
		CreatedAt:                sc.CreatedAt,
	}
}
