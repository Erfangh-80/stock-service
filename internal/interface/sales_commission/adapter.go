package salescommissioninterface

import (
	"errors"
	"time"

	app "stock-service/internal/application/sales_commission"
	domain "stock-service/internal/domain/sales_commission"
	iface "stock-service/internal/interface"
)

type SalesCommissionOutput struct {
	ID                       int64     `json:"id"`
	InventoryID              int64     `json:"inventory_id"`
	CategoryCommissionRuleID int64     `json:"category_commission_rule_id"`
	SaleModel                string    `json:"sale_model"`
	RatePercent              float64   `json:"rate_percent"`
	MinQty                   *int      `json:"min_qty,omitempty"`
	MinPrice                 float64   `json:"min_price"`
	MaxPrice                 *float64  `json:"max_price,omitempty"`
	CreatedAt                time.Time `json:"created_at"`
}

type SalesCommissionListItem struct {
	ID          int64   `json:"id"`
	InventoryID int64   `json:"inventory_id"`
	SaleModel   string  `json:"sale_model"`
	RatePercent float64 `json:"rate_percent"`
	MinPrice    float64 `json:"min_price"`
	CreatedAt   time.Time `json:"created_at"`
}

type SalesCommissionListResponse struct {
	Commissions []SalesCommissionListItem `json:"commissions"`
	Total       int                       `json:"total"`
	Page        int                       `json:"page"`
	Limit       int                       `json:"limit"`
}

type CommissionCalculationOutput struct {
	CommissionID   int64   `json:"commission_id"`
	InventoryID    int64   `json:"inventory_id"`
	BasePriceUsed  float64 `json:"base_price_used"`
	Quantity       int     `json:"quantity"`
	RatePercent    float64 `json:"rate_percent"`
	CommissionAmt  float64 `json:"commission_amt"`
	PriceSource    string  `json:"price_source"`
}

type CreateSalesCommissionParams struct {
	InventoryID              int64
	CategoryCommissionRuleID int64
	SaleModel                string
	RatePercent              float64
	MinPrice                 float64
}

type ListSalesCommissionsParams struct {
	InventoryID *int64
	SaleModel   *string
	Page        int
	Limit       int
}

type CalculateCommissionParams struct {
	InventoryID int64
	Quantity    int
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

type getSalesCommissionUseCase interface {
	Execute(input app.GetSalesCommissionInput) (*domain.SalesCommission, error)
}

type getByInventorySalesCommissionUseCase interface {
	Execute(input app.GetByInventorySalesCommissionInput) (*domain.SalesCommission, error)
}

type listSalesCommissionsUseCase interface {
	Execute(input app.ListSalesCommissionsInput) (*app.ListSalesCommissionsOutput, error)
}

type deleteSalesCommissionUseCase interface {
	Execute(input app.DeleteSalesCommissionInput) error
}

type calculateCommissionUseCase interface {
	Execute(input app.CalculateCommissionInput) (*app.CommissionCalculation, error)
}

type Adapter struct {
	create          createSalesCommissionUseCase
	updateMaxPrice  updateMaxPriceUseCase
	updateMinQty    updateMinQtyUseCase
	get             getSalesCommissionUseCase
	getByInventory  getByInventorySalesCommissionUseCase
	list            listSalesCommissionsUseCase
	deleteUC        deleteSalesCommissionUseCase
	calculate       calculateCommissionUseCase
}

func NewAdapter(
	create createSalesCommissionUseCase,
	updateMaxPrice updateMaxPriceUseCase,
	updateMinQty updateMinQtyUseCase,
	get getSalesCommissionUseCase,
	getByInventory getByInventorySalesCommissionUseCase,
	list listSalesCommissionsUseCase,
	deleteUC deleteSalesCommissionUseCase,
	calculate calculateCommissionUseCase,
) *Adapter {
	return &Adapter{
		create: create, updateMaxPrice: updateMaxPrice, updateMinQty: updateMinQty,
		get: get, getByInventory: getByInventory, list: list,
		deleteUC: deleteUC, calculate: calculate,
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

func (a *Adapter) UpdateMaxPrice(params struct {
	CommissionID int64
	MaxPrice     float64
}) error {
	err := a.updateMaxPrice.Execute(params.CommissionID, params.MaxPrice)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) UpdateMinQty(params struct {
	CommissionID int64
	MinQty       int
}) error {
	err := a.updateMinQty.Execute(params.CommissionID, params.MinQty)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Get(id int64) (*SalesCommissionOutput, error) {
	result, err := a.get.Execute(app.GetSalesCommissionInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) GetByInventory(inventoryID int64) (*SalesCommissionOutput, error) {
	result, err := a.getByInventory.Execute(app.GetByInventorySalesCommissionInput{InventoryID: inventoryID})
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) List(params ListSalesCommissionsParams) (*SalesCommissionListResponse, error) {
	input := app.ListSalesCommissionsInput{
		InventoryID: params.InventoryID,
		SaleModel:   params.SaleModel,
		Page:        params.Page,
		Limit:       params.Limit,
	}
	result, err := a.list.Execute(input)
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]SalesCommissionListItem, 0, len(result.Commissions))
	for _, sc := range result.Commissions {
		items = append(items, SalesCommissionListItem{
			ID: sc.ID, InventoryID: sc.InventoryID,
			SaleModel: string(sc.SaleModel), RatePercent: sc.RatePercent,
			MinPrice: sc.MinPrice, CreatedAt: sc.CreatedAt,
		})
	}
	return &SalesCommissionListResponse{
		Commissions: items,
		Total:       result.Total,
		Page:        result.Page,
		Limit:       result.Limit,
	}, nil
}

func (a *Adapter) Delete(id int64) error {
	err := a.deleteUC.Execute(app.DeleteSalesCommissionInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Calculate(params CalculateCommissionParams) (*CommissionCalculationOutput, error) {
	result, err := a.calculate.Execute(app.CalculateCommissionInput{
		InventoryID: params.InventoryID,
		Quantity:    params.Quantity,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &CommissionCalculationOutput{
		CommissionID:   result.CommissionID,
		InventoryID:    result.InventoryID,
		BasePriceUsed:  result.BasePriceUsed,
		Quantity:       result.Quantity,
		RatePercent:    result.RatePercent,
		CommissionAmt:  result.CommissionAmt,
		PriceSource:    result.PriceSource,
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidRatePercent),
		errors.Is(err, domain.ErrInvalidMinPrice),
		errors.Is(err, domain.ErrInvalidMaxPrice),
		errors.Is(err, domain.ErrInvalidMinQty):
		return iface.ErrInvalidInput
	case errors.Is(err, domain.ErrCommissionNotFound),
		errors.Is(err, domain.ErrRuleNotFound):
		return iface.ErrNotFound
	default:
		return iface.ErrInternal
	}
}

func toOutput(sc *domain.SalesCommission) *SalesCommissionOutput {
	return &SalesCommissionOutput{
		ID:                       sc.ID,
		InventoryID:              sc.InventoryID,
		CategoryCommissionRuleID: sc.CategoryCommissionRuleID,
		SaleModel:                string(sc.SaleModel),
		RatePercent:              sc.RatePercent,
		MinQty:                   sc.MinQty,
		MinPrice:                 sc.MinPrice,
		MaxPrice:                 sc.MaxPrice,
		CreatedAt:                sc.CreatedAt,
	}
}
