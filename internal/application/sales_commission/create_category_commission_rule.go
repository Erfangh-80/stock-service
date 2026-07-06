package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type CreateCategoryCommissionRuleInput struct {
	CategoryID  int32
	RatePercent float64
	MinPrice    float64
}

type CreateCategoryCommissionRuleUseCase struct {
	repo domainsalescommission.CategoryCommissionRuleRepository
}

func NewCreateCategoryCommissionRuleUseCase(repo domainsalescommission.CategoryCommissionRuleRepository) *CreateCategoryCommissionRuleUseCase {
	return &CreateCategoryCommissionRuleUseCase{repo: repo}
}

func (uc *CreateCategoryCommissionRuleUseCase) Execute(input CreateCategoryCommissionRuleInput) (*domainsalescommission.CategoryCommissionRule, error) {
	rule, err := domainsalescommission.NewCategoryCommissionRule(input.CategoryID, input.RatePercent, input.MinPrice)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(rule); err != nil {
		return nil, err
	}
	return rule, nil
}
