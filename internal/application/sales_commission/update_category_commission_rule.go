package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type UpdateCategoryCommissionRuleInput struct {
	ID          int64
	RatePercent *float64
	MinPrice    *float64
	MaxPrice    *float64
	Activate    *bool
}

type UpdateCategoryCommissionRuleUseCase struct {
	repo domainsalescommission.CategoryCommissionRuleRepository
}

func NewUpdateCategoryCommissionRuleUseCase(repo domainsalescommission.CategoryCommissionRuleRepository) *UpdateCategoryCommissionRuleUseCase {
	return &UpdateCategoryCommissionRuleUseCase{repo: repo}
}

func (uc *UpdateCategoryCommissionRuleUseCase) Execute(input UpdateCategoryCommissionRuleInput) (*domainsalescommission.CategoryCommissionRule, error) {
	rule, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, domainsalescommission.ErrRuleNotFound
	}

	if input.RatePercent != nil {
		if err := domainsalescommission.ValidateRatePercent(*input.RatePercent); err != nil {
			return nil, err
		}
		rule.RatePercent = *input.RatePercent
	}
	if input.MinPrice != nil {
		if err := domainsalescommission.ValidateMinPrice(*input.MinPrice); err != nil {
			return nil, err
		}
		rule.MinPrice = *input.MinPrice
	}
	if input.MaxPrice != nil {
		if err := rule.UpdateMaxPrice(*input.MaxPrice); err != nil {
			return nil, err
		}
	}
	if input.Activate != nil {
		if *input.Activate {
			rule.Activate()
		} else {
			rule.Deactivate()
		}
	}

	if err := uc.repo.Save(rule); err != nil {
		return nil, err
	}
	return rule, nil
}
