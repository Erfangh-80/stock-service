package salescommissioninterface

import (
	"time"

	app "stock-service/internal/application/sales_commission"
	domain "stock-service/internal/domain/sales_commission"
)

type CategoryCommissionRuleOutput struct {
	ID          int64     `json:"id"`
	CategoryID  int32     `json:"category_id"`
	RatePercent float64   `json:"rate_percent"`
	MinPrice    float64   `json:"min_price"`
	MaxPrice    *float64  `json:"max_price,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type CategoryCommissionRuleListItem struct {
	ID          int64   `json:"id"`
	CategoryID  int32   `json:"category_id"`
	RatePercent float64 `json:"rate_percent"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type CategoryCommissionRuleListResponse struct {
	Rules []CategoryCommissionRuleListItem `json:"rules"`
	Total int                              `json:"total"`
	Page  int                              `json:"page"`
	Limit int                              `json:"limit"`
}

type CreateCategoryCommissionRuleParams struct {
	CategoryID  int32   `json:"category_id"`
	RatePercent float64 `json:"rate_percent"`
	MinPrice    float64 `json:"min_price"`
}

type UpdateCategoryCommissionRuleParams struct {
	RatePercent *float64 `json:"rate_percent,omitempty"`
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
	Activate    *bool    `json:"activate,omitempty"`
}

type ListCategoryCommissionRulesParams struct {
	CategoryID *int32
	IsActive   *bool
	Page       int
	Limit      int
}

type createCategoryCommissionRuleUseCase interface {
	Execute(input app.CreateCategoryCommissionRuleInput) (*domain.CategoryCommissionRule, error)
}

type getCategoryCommissionRuleUseCase interface {
	Execute(input app.GetCategoryCommissionRuleInput) (*domain.CategoryCommissionRule, error)
}

type listCategoryCommissionRulesUseCase interface {
	Execute(input app.ListCategoryCommissionRulesInput) (*app.ListCategoryCommissionRulesOutput, error)
}

type updateCategoryCommissionRuleUseCase interface {
	Execute(input app.UpdateCategoryCommissionRuleInput) (*domain.CategoryCommissionRule, error)
}

type deleteCategoryCommissionRuleUseCase interface {
	Execute(input app.DeleteCategoryCommissionRuleInput) error
}

type CategoryCommissionRuleAdapter struct {
	create   createCategoryCommissionRuleUseCase
	get      getCategoryCommissionRuleUseCase
	list     listCategoryCommissionRulesUseCase
	update   updateCategoryCommissionRuleUseCase
	deleteUC deleteCategoryCommissionRuleUseCase
}

func NewCategoryCommissionRuleAdapter(
	create createCategoryCommissionRuleUseCase,
	get getCategoryCommissionRuleUseCase,
	list listCategoryCommissionRulesUseCase,
	update updateCategoryCommissionRuleUseCase,
	deleteUC deleteCategoryCommissionRuleUseCase,
) *CategoryCommissionRuleAdapter {
	return &CategoryCommissionRuleAdapter{
		create: create, get: get, list: list,
		update: update, deleteUC: deleteUC,
	}
}

func (a *CategoryCommissionRuleAdapter) Create(params CreateCategoryCommissionRuleParams) (*CategoryCommissionRuleOutput, error) {
	result, err := a.create.Execute(app.CreateCategoryCommissionRuleInput{
		CategoryID: params.CategoryID, RatePercent: params.RatePercent, MinPrice: params.MinPrice,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toRuleOutput(result), nil
}

func (a *CategoryCommissionRuleAdapter) Get(id int64) (*CategoryCommissionRuleOutput, error) {
	result, err := a.get.Execute(app.GetCategoryCommissionRuleInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toRuleOutput(result), nil
}

func (a *CategoryCommissionRuleAdapter) List(params ListCategoryCommissionRulesParams) (*CategoryCommissionRuleListResponse, error) {
	input := app.ListCategoryCommissionRulesInput{
		CategoryID: params.CategoryID,
		IsActive:   params.IsActive,
		Page:       params.Page,
		Limit:      params.Limit,
	}
	result, err := a.list.Execute(input)
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]CategoryCommissionRuleListItem, 0, len(result.Rules))
	for _, r := range result.Rules {
		items = append(items, CategoryCommissionRuleListItem{
			ID: r.ID, CategoryID: r.CategoryID,
			RatePercent: r.RatePercent, IsActive: r.IsActive,
			CreatedAt: r.CreatedAt,
		})
	}
	return &CategoryCommissionRuleListResponse{
		Rules: items, Total: result.Total,
		Page: result.Page, Limit: result.Limit,
	}, nil
}

func (a *CategoryCommissionRuleAdapter) Update(id int64, params UpdateCategoryCommissionRuleParams) (*CategoryCommissionRuleOutput, error) {
	result, err := a.update.Execute(app.UpdateCategoryCommissionRuleInput{
		ID: id, RatePercent: params.RatePercent, MinPrice: params.MinPrice,
		MaxPrice: params.MaxPrice, Activate: params.Activate,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toRuleOutput(result), nil
}

func (a *CategoryCommissionRuleAdapter) Delete(id int64) error {
	err := a.deleteUC.Execute(app.DeleteCategoryCommissionRuleInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func toRuleOutput(r *domain.CategoryCommissionRule) *CategoryCommissionRuleOutput {
	return &CategoryCommissionRuleOutput{
		ID: r.ID, CategoryID: r.CategoryID,
		RatePercent: r.RatePercent, MinPrice: r.MinPrice,
		MaxPrice: r.MaxPrice, IsActive: r.IsActive,
		CreatedAt: r.CreatedAt,
	}
}
