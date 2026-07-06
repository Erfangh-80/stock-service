package productbundleinterface

import (
	appproductbundle "stock-service/internal/application/product_bundle"
	"stock-service/internal/domain/product_bundle"
	iface "stock-service/internal/interface"
)

type ProductBundleResponse struct {
	ID               int64  `json:"id"`
	ProductID        int32  `json:"product_id"`
	RelatedProductID int32  `json:"related_product_id"`
	Type             string `json:"type"`
	SortOrder        int    `json:"sort_order"`
	CreatedAt        string `json:"created_at"`
}

type CreateBundleParams struct {
	ProductID        int32
	RelatedProductID int32
	Type             string
	SortOrder        int
}

type Adapter struct {
	create appproductbundle.CreateBundleUseCase
	list   appproductbundle.ListBundlesUseCase
}

func NewAdapter(
	create appproductbundle.CreateBundleUseCase,
	list appproductbundle.ListBundlesUseCase,
) *Adapter {
	return &Adapter{create: create, list: list}
}

func (a *Adapter) Create(params CreateBundleParams) (*ProductBundleResponse, error) {
	result, err := a.create.Execute(appproductbundle.CreateBundleInput{
		ProductID: params.ProductID, RelatedProductID: params.RelatedProductID,
		Type: params.Type, SortOrder: params.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) List(productID int32) ([]*ProductBundleResponse, error) {
	result, err := a.list.Execute(appproductbundle.ListBundlesInput{ProductID: productID})
	if err != nil {
		return nil, mapError(err)
	}
	resp := make([]*ProductBundleResponse, len(result.Bundles))
	for i, pb := range result.Bundles {
		resp[i] = toResponse(pb)
	}
	return resp, nil
}

func mapError(err error) error {
	switch err {
	case productbundle.ErrBundleNotFound:
		return iface.ErrNotFound
	case productbundle.ErrInvalidProductID, productbundle.ErrInvalidRelatedProductID,
		productbundle.ErrSelfReference, productbundle.ErrInvalidBundleType:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(pb *productbundle.ProductBundle) *ProductBundleResponse {
	return &ProductBundleResponse{
		ID: pb.ID, ProductID: pb.ProductID,
		RelatedProductID: pb.RelatedProductID, Type: string(pb.Type),
		SortOrder: pb.SortOrder, CreatedAt: pb.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
