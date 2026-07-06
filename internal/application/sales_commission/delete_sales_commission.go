package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type DeleteSalesCommissionInput struct {
	ID int64
}

type DeleteSalesCommissionUseCase struct {
	repo domainsalescommission.Repository
}

func NewDeleteSalesCommissionUseCase(repo domainsalescommission.Repository) *DeleteSalesCommissionUseCase {
	return &DeleteSalesCommissionUseCase{repo: repo}
}

func (uc *DeleteSalesCommissionUseCase) Execute(input DeleteSalesCommissionInput) error {
	sc, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if sc == nil {
		return domainsalescommission.ErrCommissionNotFound
	}
	return uc.repo.Delete(input.ID)
}
