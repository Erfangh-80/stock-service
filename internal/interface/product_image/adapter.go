package productimageinterface

import (
	appproductimage "stock-service/internal/application/product_image"
	"stock-service/internal/domain/product_image"
	iface "stock-service/internal/interface"
)

type ProductImageResponse struct {
	ID        int64  `json:"id"`
	ProductID int32  `json:"product_id"`
	FileID    int64  `json:"file_id"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

type CreateImageParams struct {
	ProductID int32
	FileID    int64
	SortOrder int
}

type Adapter struct {
	create appproductimage.CreateImageUseCase
	list   appproductimage.ListImagesUseCase
	delete appproductimage.DeleteImageUseCase
}

func NewAdapter(
	create appproductimage.CreateImageUseCase,
	list appproductimage.ListImagesUseCase,
	delete appproductimage.DeleteImageUseCase,
) *Adapter {
	return &Adapter{create: create, list: list, delete: delete}
}

func (a *Adapter) Create(params CreateImageParams) (*ProductImageResponse, error) {
	result, err := a.create.Execute(appproductimage.CreateImageInput{
		ProductID: params.ProductID, FileID: params.FileID, SortOrder: params.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) List(productID int32) ([]*ProductImageResponse, error) {
	result, err := a.list.Execute(appproductimage.ListImagesInput{ProductID: productID})
	if err != nil {
		return nil, mapError(err)
	}
	resp := make([]*ProductImageResponse, len(result.Images))
	for i, img := range result.Images {
		resp[i] = toResponse(img)
	}
	return resp, nil
}

func (a *Adapter) Delete(id int64) error {
	err := a.delete.Execute(appproductimage.DeleteImageInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func mapError(err error) error {
	switch err {
	case productimage.ErrImageNotFound:
		return iface.ErrNotFound
	case productimage.ErrInvalidProductID, productimage.ErrInvalidFileID:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(img *productimage.ProductImage) *ProductImageResponse {
	return &ProductImageResponse{
		ID: img.ID, ProductID: img.ProductID, FileID: img.FileID,
		SortOrder: img.SortOrder, CreatedAt: img.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
