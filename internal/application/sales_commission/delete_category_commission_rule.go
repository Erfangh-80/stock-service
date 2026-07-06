package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type DeleteCategoryCommissionRuleInput struct {
	ID int64
}

type DeleteCategoryCommissionRuleUseCase struct {
	repo domainsalescommission.CategoryCommissionRuleRepository
}

func NewDeleteCategoryCommissionRuleUseCase(repo domainsalescommission.CategoryCommissionRuleRepository) *DeleteCategoryCommissionRuleUseCase {
	return &DeleteCategoryCommissionRuleUseCase{repo: repo}
}

func (uc *DeleteCategoryCommissionRuleUseCase) Execute(input DeleteCategoryCommissionRuleInput) error {
	rule, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if rule == nil {
		return domainsalescommission.ErrRuleNotFound
	}
	return uc.repo.Delete(input.ID)
}
