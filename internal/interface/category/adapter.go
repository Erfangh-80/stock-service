package categoryinterface

import (
	appcategory "stock-service/internal/application/category"
	"stock-service/internal/domain/category"
	iface "stock-service/internal/interface"
)

type CategoryResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	ParentID    *int64  `json:"parent_id,omitempty"`
	Description *string `json:"description,omitempty"`
	ImageFileID *int64  `json:"image_file_id,omitempty"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateCategoryParams struct {
	Name        string
	Slug        string
	ParentID    *int64
	Description *string
}

type UpdateCategoryParams struct {
	ID          int64
	Name        *string
	Slug        *string
	Description *string
	ParentID    *int64
	SortOrder   *int
}

type Adapter struct {
	create appcategory.CreateCategoryUseCase
	get    appcategory.GetCategoryUseCase
	update appcategory.UpdateCategoryUseCase
	delete appcategory.DeleteCategoryUseCase
	list   appcategory.ListCategoriesUseCase
}

func NewAdapter(
	create appcategory.CreateCategoryUseCase,
	get appcategory.GetCategoryUseCase,
	update appcategory.UpdateCategoryUseCase,
	delete appcategory.DeleteCategoryUseCase,
	list appcategory.ListCategoriesUseCase,
) *Adapter {
	return &Adapter{
		create: create, get: get, update: update,
		delete: delete, list: list,
	}
}

func (a *Adapter) Create(params CreateCategoryParams) (*CategoryResponse, error) {
	result, err := a.create.Execute(appcategory.CreateCategoryInput{
		Name: params.Name, Slug: params.Slug,
		ParentID: params.ParentID, Description: params.Description,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Get(id int64) (*CategoryResponse, error) {
	result, err := a.get.Execute(appcategory.GetCategoryInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Update(params UpdateCategoryParams) (*CategoryResponse, error) {
	result, err := a.update.Execute(appcategory.UpdateCategoryInput{
		ID: params.ID, Name: params.Name, Slug: params.Slug,
		Description: params.Description, ParentID: params.ParentID,
		SortOrder: params.SortOrder,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Delete(id int64) error {
	err := a.delete.Execute(appcategory.DeleteCategoryInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) List() ([]*CategoryResponse, error) {
	result, err := a.list.Execute()
	if err != nil {
		return nil, mapError(err)
	}
	resp := make([]*CategoryResponse, len(result.Categories))
	for i, c := range result.Categories {
		resp[i] = toResponse(c)
	}
	return resp, nil
}

func mapError(err error) error {
	switch err {
	case category.ErrCategoryNotFound:
		return iface.ErrNotFound
	case category.ErrNameRequired, category.ErrSlugRequired:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(c *category.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		ParentID:    c.ParentID,
		Description: c.Description,
		ImageFileID: c.ImageFileID,
		SortOrder:   c.SortOrder,
		CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
