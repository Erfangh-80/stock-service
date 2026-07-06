package brandinterface

import (
	appbrand "stock-service/internal/application/brand"
	"stock-service/internal/domain/brand"
	iface "stock-service/internal/interface"
)

type BrandResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	LogoFileID *int64 `json:"logo_file_id,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateBrandParams struct {
	Name string
	Slug string
}

type UpdateBrandParams struct {
	ID   int64
	Name *string
	Slug *string
}

type Adapter struct {
	create appbrand.CreateBrandUseCase
	get    appbrand.GetBrandUseCase
	update appbrand.UpdateBrandUseCase
	delete appbrand.DeleteBrandUseCase
	list   appbrand.ListBrandsUseCase
}

func NewAdapter(
	create appbrand.CreateBrandUseCase,
	get appbrand.GetBrandUseCase,
	update appbrand.UpdateBrandUseCase,
	delete appbrand.DeleteBrandUseCase,
	list appbrand.ListBrandsUseCase,
) *Adapter {
	return &Adapter{
		create: create, get: get, update: update,
		delete: delete, list: list,
	}
}

func (a *Adapter) Create(params CreateBrandParams) (*BrandResponse, error) {
	result, err := a.create.Execute(appbrand.CreateBrandInput{
		Name: params.Name, Slug: params.Slug,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Get(id int64) (*BrandResponse, error) {
	result, err := a.get.Execute(appbrand.GetBrandInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Update(params UpdateBrandParams) (*BrandResponse, error) {
	result, err := a.update.Execute(appbrand.UpdateBrandInput{
		ID: params.ID, Name: params.Name, Slug: params.Slug,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Delete(id int64) error {
	err := a.delete.Execute(appbrand.DeleteBrandInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) List() ([]*BrandResponse, error) {
	result, err := a.list.Execute()
	if err != nil {
		return nil, mapError(err)
	}
	resp := make([]*BrandResponse, len(result.Brands))
	for i, b := range result.Brands {
		resp[i] = toResponse(b)
	}
	return resp, nil
}

func mapError(err error) error {
	switch err {
	case brand.ErrBrandNotFound:
		return iface.ErrNotFound
	case brand.ErrNameRequired, brand.ErrSlugRequired:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(b *brand.Brand) *BrandResponse {
	return &BrandResponse{
		ID:         b.ID,
		Name:       b.Name,
		Slug:       b.Slug,
		LogoFileID: b.LogoFileID,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
