package referencepriceinterface

import (
	"errors"
	"time"

	domain "stock-service/internal/domain/reference_price"
	iface "stock-service/internal/interface"
)

type ReferencePriceOutput struct {
	ID        int64
	ProductID int32
	Price     float64
	Source    string
	CreatedAt time.Time
}

type CreateReferencePriceParams struct {
	ProductID int32
	Price     float64
	Source    string
}

type createReferencePriceUseCase interface {
	Execute(productID int32, price float64, source string) (*domain.ReferencePrice, error)
}

type Adapter struct {
	create createReferencePriceUseCase
}

func NewAdapter(create createReferencePriceUseCase) *Adapter {
	return &Adapter{create: create}
}

func (a *Adapter) Create(params CreateReferencePriceParams) (*ReferencePriceOutput, error) {
	result, err := a.create.Execute(params.ProductID, params.Price, params.Source)
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidReferencePrice):
		return iface.ErrInvalidInput
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
