package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type ListCategoryCommissionRulesInput struct {
	CategoryID *int32
	IsActive   *bool
	Page       int
	Limit      int
}

type ListCategoryCommissionRulesOutput struct {
	Rules []*domainsalescommission.CategoryCommissionRule
	Total int
	Page  int
	Limit int
}

type ListCategoryCommissionRulesUseCase struct {
	repo domainsalescommission.CategoryCommissionRuleRepository
}

func NewListCategoryCommissionRulesUseCase(repo domainsalescommission.CategoryCommissionRuleRepository) *ListCategoryCommissionRulesUseCase {
	return &ListCategoryCommissionRulesUseCase{repo: repo}
}

func (uc *ListCategoryCommissionRulesUseCase) Execute(input ListCategoryCommissionRulesInput) (*ListCategoryCommissionRulesOutput, error) {
	filter := domainsalescommission.CategoryCommissionRuleFilter{
		CategoryID: input.CategoryID,
		IsActive:   input.IsActive,
		Page:       input.Page,
		Limit:      input.Limit,
	}
	items, total, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return &ListCategoryCommissionRulesOutput{
		Rules: items,
		Total: total,
		Page:  input.Page,
		Limit: input.Limit,
	}, nil
}
