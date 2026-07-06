package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type GetCategoryCommissionRuleInput struct {
	ID int64
}

type GetCategoryCommissionRuleUseCase struct {
	repo domainsalescommission.CategoryCommissionRuleRepository
}

func NewGetCategoryCommissionRuleUseCase(repo domainsalescommission.CategoryCommissionRuleRepository) *GetCategoryCommissionRuleUseCase {
	return &GetCategoryCommissionRuleUseCase{repo: repo}
}

func (uc *GetCategoryCommissionRuleUseCase) Execute(input GetCategoryCommissionRuleInput) (*domainsalescommission.CategoryCommissionRule, error) {
	rule, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, domainsalescommission.ErrRuleNotFound
	}
	return rule, nil
}
