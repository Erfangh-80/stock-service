package referencepriceinterface

import (
	"errors"
	"time"

	app "stock-service/internal/application/reference_price"
	domain "stock-service/internal/domain/reference_price"
	iface "stock-service/internal/interface"
)

type ReferencePriceOutput struct {
	ID        int64     `json:"id"`
	ProductID int32     `json:"product_id"`
	Price     float64   `json:"price"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

type ReferencePriceListItem struct {
	ID        int64     `json:"id"`
	ProductID int32     `json:"product_id"`
	Price     float64   `json:"price"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

type ReferencePriceListResponse struct {
	ReferencePrices []ReferencePriceListItem `json:"reference_prices"`
	Total           int                      `json:"total"`
	Page            int                      `json:"page"`
	Limit           int                      `json:"limit"`
}

type ValidationOutput struct {
	ProductID        int32     `json:"product_id"`
	ReferencePriceID int64     `json:"reference_price_id"`
	ReferencePrice   float64   `json:"reference_price"`
	Source           string    `json:"source"`
	InventoryCount   int       `json:"inventory_count"`
	BasePrices       []float64 `json:"base_prices"`
	Comparison       string    `json:"comparison"`
}

type CreateReferencePriceParams struct {
	ProductID int32   `json:"product_id"`
	Price     float64 `json:"price"`
	Source    string  `json:"source"`
}

type ListReferencePricesParams struct {
	ProductID *int32
	Source    *string
	Page      int
	Limit     int
}

type createReferencePriceUseCase interface {
	Execute(productID int32, price float64, source string) (*domain.ReferencePrice, error)
}

type getReferencePriceUseCase interface {
	Execute(input app.GetReferencePriceInput) (*domain.ReferencePrice, error)
}

type getByProductReferencePriceUseCase interface {
	Execute(input app.GetByProductReferencePriceInput) (*domain.ReferencePrice, error)
}

type listReferencePricesUseCase interface {
	Execute(input app.ListReferencePricesInput) (*app.ListReferencePricesOutput, error)
}

type deleteReferencePriceUseCase interface {
	Execute(input app.DeleteReferencePriceInput) error
}

type validateReferencePriceUseCase interface {
	Execute(input app.ValidateReferencePriceInput) (*app.ReferencePriceValidation, error)
}

type Adapter struct {
	create       createReferencePriceUseCase
	get          getReferencePriceUseCase
	getByProduct  getByProductReferencePriceUseCase
	list         listReferencePricesUseCase
	deleteUC     deleteReferencePriceUseCase
	validate     validateReferencePriceUseCase
}

func NewAdapter(
	create createReferencePriceUseCase,
	get getReferencePriceUseCase,
	getByProduct getByProductReferencePriceUseCase,
	list listReferencePricesUseCase,
	deleteUC deleteReferencePriceUseCase,
	validate validateReferencePriceUseCase,
) *Adapter {
	return &Adapter{
		create: create, get: get, getByProduct: getByProduct,
		list: list, deleteUC: deleteUC, validate: validate,
	}
}

func (a *Adapter) Create(params CreateReferencePriceParams) (*ReferencePriceOutput, error) {
	result, err := a.create.Execute(params.ProductID, params.Price, params.Source)
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) Get(id int64) (*ReferencePriceOutput, error) {
	result, err := a.get.Execute(app.GetReferencePriceInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) GetByProduct(productID int32) (*ReferencePriceOutput, error) {
	result, err := a.getByProduct.Execute(app.GetByProductReferencePriceInput{ProductID: productID})
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) List(params ListReferencePricesParams) (*ReferencePriceListResponse, error) {
	input := app.ListReferencePricesInput{
		ProductID: params.ProductID,
		Source:    params.Source,
		Page:      params.Page,
		Limit:     params.Limit,
	}
	result, err := a.list.Execute(input)
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]ReferencePriceListItem, 0, len(result.ReferencePrices))
	for _, rp := range result.ReferencePrices {
		items = append(items, ReferencePriceListItem{
			ID: rp.ID, ProductID: rp.ProductID,
			Price: rp.Price, Source: rp.Source,
			CreatedAt: rp.CreatedAt,
		})
	}
	return &ReferencePriceListResponse{
		ReferencePrices: items,
		Total:           result.Total,
		Page:            result.Page,
		Limit:           result.Limit,
	}, nil
}

func (a *Adapter) Delete(id int64) error {
	err := a.deleteUC.Execute(app.DeleteReferencePriceInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Validate(productID int32) (*ValidationOutput, error) {
	result, err := a.validate.Execute(app.ValidateReferencePriceInput{ProductID: productID})
	if err != nil {
		return nil, mapError(err)
	}
	return &ValidationOutput{
		ProductID:        result.ProductID,
		ReferencePriceID: result.ReferencePriceID,
		ReferencePrice:   result.ReferencePrice,
		Source:           result.Source,
		InventoryCount:   result.InventoryCount,
		BasePrices:       result.BasePrices,
		Comparison:       result.Comparison,
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidReferencePrice):
		return iface.ErrInvalidInput
	case errors.Is(err, domain.ErrReferencePriceNotFound):
		return iface.ErrNotFound
	default:
		return iface.ErrInternal
	}
}

func toOutput(rp *domain.ReferencePrice) *ReferencePriceOutput {
	return &ReferencePriceOutput{
		ID:        rp.ID,
		ProductID: rp.ProductID,
		Price:     rp.Price,
		Source:    rp.Source,
		CreatedAt: rp.CreatedAt,
	}
}
