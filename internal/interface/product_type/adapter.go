package producttypeinterface

import (
	appproducttype "stock-service/internal/application/product_type"
	"stock-service/internal/domain/product_type"
	iface "stock-service/internal/interface"
)

type ProductTypeResponse struct {
	ID        int64  `json:"id"`
	ProductID int32  `json:"product_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

type CreateTypeParams struct {
	ProductID int32
	Name      string
	Value     string
	SortOrder int
}

type Adapter struct {
	create appproducttype.CreateTypeUseCase
	list   appproducttype.ListTypesUseCase
}

func NewAdapter(
	create appproducttype.CreateTypeUseCase,
	list appproducttype.ListTypesUseCase,
) *Adapter {
	return &Adapter{create: create, list: list}
}

func (a *Adapter) Create(params CreateTypeParams) (*ProductTypeResponse, error) {
	result, err := a.create.Execute(appproducttype.CreateTypeInput{
		ProductID: params.ProductID, Name: params.Name,
		Value: params.Value, SortOrder: params.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) List(productID int32) ([]*ProductTypeResponse, error) {
	result, err := a.list.Execute(appproducttype.ListTypesInput{ProductID: productID})
	if err != nil {
		return nil, mapError(err)
	}
	resp := make([]*ProductTypeResponse, len(result.Types))
	for i, pt := range result.Types {
		resp[i] = toResponse(pt)
	}
	return resp, nil
}

func mapError(err error) error {
	switch err {
	case producttype.ErrTypeNotFound:
		return iface.ErrNotFound
	case producttype.ErrInvalidProductID, producttype.ErrNameRequired, producttype.ErrValueRequired:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(pt *producttype.ProductType) *ProductTypeResponse {
	return &ProductTypeResponse{
		ID: pt.ID, ProductID: pt.ProductID, Name: pt.Name,
		Value: pt.Value, SortOrder: pt.SortOrder,
		CreatedAt: pt.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
