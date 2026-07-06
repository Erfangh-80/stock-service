package productinterface

import (
	appproduct "stock-service/internal/application/product"
	"stock-service/internal/domain/product"
	iface "stock-service/internal/interface"
	"time"
)

type ProductResponse struct {
	ID               int32   `json:"id"`
	TitleFa          string  `json:"title_fa"`
	TitleEn          *string `json:"title_en,omitempty"`
	Slug             string  `json:"slug,omitempty"`
	Description      *string `json:"description,omitempty"`
	BrandID          int64   `json:"brand_id"`
	CategoryID       int64   `json:"category_id"`
	OwnerType        string  `json:"owner_type"`
	OwnerID          *int64  `json:"owner_id,omitempty"`
	IsOriginal       bool    `json:"is_original"`
	IsEnabled        bool    `json:"is_enabled"`
	EnabledAt        *string `json:"enabled_at,omitempty"`
	DisabledAt       *string `json:"disabled_at,omitempty"`
	MetaTitle        *string `json:"meta_title,omitempty"`
	MetaDescription  *string `json:"meta_description,omitempty"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	IndexImageFileID *int64  `json:"index_image_file_id,omitempty"`
	DeletedAt        *string `json:"deleted_at,omitempty"`
}

type ProductListResponse struct {
	Products []*ProductResponse `json:"products"`
	Total    int                `json:"total"`
	Page     int                `json:"page"`
	Limit    int                `json:"limit"`
}

type CreateProductParams struct {
	TitleFa          string
	TitleEn          *string
	Slug             string
	Description      *string
	BrandID          int64
	CategoryID       int64
	OwnerType        string
	OwnerID          *int64
	IsOriginal       *bool
	MetaTitle        *string
	MetaDescription  *string
	IndexImageFileID *int64
}

type UpdateProductParams struct {
	ID               int32
	TitleFa          *string
	TitleEn          *string
	Slug             *string
	Description      *string
	BrandID          *int64
	CategoryID       *int64
	MetaTitle        *string
	MetaDescription  *string
	IndexImageFileID *int64
}

type ListProductFilter struct {
	OwnerType  *string
	OwnerID    *int64
	Status     *string
	CategoryID *int64
	BrandID    *int64
	Search     *string
	Page       int
	Limit      int
}

type createProductUseCase interface {
	Execute(input appproduct.CreateProductInput) (*product.Product, error)
}

type getProductUseCase interface {
	Execute(input appproduct.GetProductInput) (*product.Product, error)
}

type updateProductUseCase interface {
	Execute(input appproduct.UpdateProductInput) (*product.Product, error)
}

type activateProductUseCase interface {
	Execute(input appproduct.ActivateProductInput) (*product.Product, error)
}

type rejectProductUseCase interface {
	Execute(input appproduct.RejectProductInput) (*product.Product, error)
}

type softDeleteProductUseCase interface {
	Execute(input appproduct.SoftDeleteProductInput) (*product.Product, error)
}

type enableProductUseCase interface {
	Execute(input appproduct.EnableProductInput) (*product.Product, error)
}

type disableProductUseCase interface {
	Execute(input appproduct.DisableProductInput) (*product.Product, error)
}

type updateSEOUseCase interface {
	Execute(input appproduct.UpdateSEOInput) (*product.Product, error)
}

type listProductsUseCase interface {
	Execute(input appproduct.ListProductsInput) (*appproduct.ListProductsOutput, error)
}

type Adapter struct {
	create     createProductUseCase
	get        getProductUseCase
	update     updateProductUseCase
	activate   activateProductUseCase
	reject     rejectProductUseCase
	softDelete softDeleteProductUseCase
	enable     enableProductUseCase
	disable    disableProductUseCase
	updateSEO  updateSEOUseCase
	list       listProductsUseCase
}

func NewAdapter(
	create createProductUseCase,
	get getProductUseCase,
	update updateProductUseCase,
	activate activateProductUseCase,
	reject rejectProductUseCase,
	softDelete softDeleteProductUseCase,
	enable enableProductUseCase,
	disable disableProductUseCase,
	updateSEO updateSEOUseCase,
	list listProductsUseCase,
) *Adapter {
	return &Adapter{
		create: create, get: get, update: update,
		activate: activate, reject: reject, softDelete: softDelete,
		enable: enable, disable: disable, updateSEO: updateSEO,
		list: list,
	}
}

func (a *Adapter) Create(params CreateProductParams) (*ProductResponse, error) {
	var ownerType product.OwnerType
	if params.OwnerType != "" {
		ownerType = product.OwnerType(params.OwnerType)
	}

	result, err := a.create.Execute(appproduct.CreateProductInput{
		TitleFa:          params.TitleFa,
		TitleEn:          params.TitleEn,
		Slug:             params.Slug,
		Description:      params.Description,
		BrandID:          params.BrandID,
		CategoryID:       params.CategoryID,
		OwnerType:        ownerType,
		OwnerID:          params.OwnerID,
		IsOriginal:       params.IsOriginal,
		MetaTitle:        params.MetaTitle,
		MetaDescription:  params.MetaDescription,
		IndexImageFileID: params.IndexImageFileID,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Get(id int32) (*ProductResponse, error) {
	result, err := a.get.Execute(appproduct.GetProductInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Update(params UpdateProductParams) (*ProductResponse, error) {
	result, err := a.update.Execute(appproduct.UpdateProductInput{
		ID:               params.ID,
		TitleFa:          params.TitleFa,
		TitleEn:          params.TitleEn,
		Slug:             params.Slug,
		Description:      params.Description,
		BrandID:          params.BrandID,
		CategoryID:       params.CategoryID,
		MetaTitle:        params.MetaTitle,
		MetaDescription:  params.MetaDescription,
		IndexImageFileID: params.IndexImageFileID,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Activate(id int32) (*ProductResponse, error) {
	result, err := a.activate.Execute(appproduct.ActivateProductInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Reject(id int32) (*ProductResponse, error) {
	result, err := a.reject.Execute(appproduct.RejectProductInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) SoftDelete(id int32) (*ProductResponse, error) {
	result, err := a.softDelete.Execute(appproduct.SoftDeleteProductInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Enable(id int32) (*ProductResponse, error) {
	result, err := a.enable.Execute(appproduct.EnableProductInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Disable(id int32) (*ProductResponse, error) {
	result, err := a.disable.Execute(appproduct.DisableProductInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) UpdateSEO(id int32, metaTitle, metaDescription *string) (*ProductResponse, error) {
	result, err := a.updateSEO.Execute(appproduct.UpdateSEOInput{
		ID: id, MetaTitle: metaTitle, MetaDescription: metaDescription,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) List(filter ListProductFilter) (*ProductListResponse, error) {
	var ownerType *product.OwnerType
	if filter.OwnerType != nil {
		t := product.OwnerType(*filter.OwnerType)
		ownerType = &t
	}

	var status *product.ProductStatus
	if filter.Status != nil {
		t := product.ProductStatus(*filter.Status)
		status = &t
	}

	result, err := a.list.Execute(appproduct.ListProductsInput{
		OwnerType:  ownerType,
		OwnerID:    filter.OwnerID,
		Status:     status,
		CategoryID: filter.CategoryID,
		BrandID:    filter.BrandID,
		Search:     filter.Search,
		Page:       filter.Page,
		Limit:      filter.Limit,
	})
	if err != nil {
		return nil, mapError(err)
	}

	products := make([]*ProductResponse, len(result.Products))
	for i, p := range result.Products {
		products[i] = toResponse(p)
	}

	return &ProductListResponse{
		Products: products,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
	}, nil
}

func mapError(err error) error {
	switch err {
	case product.ErrProductNotFound:
		return iface.ErrNotFound
	case product.ErrTitleFaRequired, product.ErrInvalidBrandID, product.ErrInvalidCategoryID, product.ErrSlugRequired:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z")
}

func toResponse(p *product.Product) *ProductResponse {
	r := &ProductResponse{
		ID:               p.ID,
		TitleFa:          p.TitleFa,
		TitleEn:          p.TitleEn,
		Slug:             p.Slug,
		Description:      p.Description,
		BrandID:          p.BrandID,
		CategoryID:       p.CategoryID,
		OwnerType:        string(p.OwnerType),
		OwnerID:          p.OwnerID,
		IsOriginal:       p.IsOriginal,
		IsEnabled:        p.IsEnabled,
		MetaTitle:        p.MetaTitle,
		MetaDescription:  p.MetaDescription,
		Status:           string(p.Status),
		CreatedAt:        formatTime(p.CreatedAt),
		UpdatedAt:        formatTime(p.UpdatedAt),
		IndexImageFileID: p.IndexImageFileID,
	}
	if p.DeletedAt != nil {
		s := formatTime(*p.DeletedAt)
		r.DeletedAt = &s
	}
	if p.EnabledAt != nil {
		s := formatTime(*p.EnabledAt)
		r.EnabledAt = &s
	}
	if p.DisabledAt != nil {
		s := formatTime(*p.DisabledAt)
		r.DisabledAt = &s
	}
	return r
}
