package productattributeinterface

import (
	appproductattribute "stock-service/internal/application/product_attribute"
	"stock-service/internal/domain/product_attribute"
	iface "stock-service/internal/interface"
)

type ProductAttributeResponse struct {
	ID        int64  `json:"id"`
	ProductID int32  `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateAttributeParams struct {
	ProductID int32
	Key       string
	Value     string
}

type Adapter struct {
	create appproductattribute.CreateAttributeUseCase
	list   appproductattribute.ListAttributesUseCase
}

func NewAdapter(
	create appproductattribute.CreateAttributeUseCase,
	list appproductattribute.ListAttributesUseCase,
) *Adapter {
	return &Adapter{create: create, list: list}
}

func (a *Adapter) Create(params CreateAttributeParams) (*ProductAttributeResponse, error) {
	result, err := a.create.Execute(appproductattribute.CreateAttributeInput{
		ProductID: params.ProductID, Key: params.Key, Value: params.Value,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) List(productID int32) ([]*ProductAttributeResponse, error) {
	result, err := a.list.Execute(appproductattribute.ListAttributesInput{ProductID: productID})
	if err != nil {
		return nil, mapError(err)
	}
	resp := make([]*ProductAttributeResponse, len(result.Attributes))
	for i, attr := range result.Attributes {
		resp[i] = toResponse(attr)
	}
	return resp, nil
}

func mapError(err error) error {
	switch err {
	case productattribute.ErrAttributeNotFound:
		return iface.ErrNotFound
	case productattribute.ErrInvalidProductID, productattribute.ErrKeyRequired:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(pa *productattribute.ProductAttribute) *ProductAttributeResponse {
	return &ProductAttributeResponse{
		ID: pa.ID, ProductID: pa.ProductID, Key: pa.Key,
		Value: pa.Value, CreatedAt: pa.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: pa.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
