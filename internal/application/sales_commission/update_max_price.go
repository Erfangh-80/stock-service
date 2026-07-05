package salescommission

import (
	"stock-service/internal/domain/sales_commission"
)

type UpdateMaxPriceUseCase struct {
	repo salescommission.Repository
}

func NewUpdateMaxPriceUseCase(repo salescommission.Repository) *UpdateMaxPriceUseCase {
	return &UpdateMaxPriceUseCase{repo: repo}
}

func (uc *UpdateMaxPriceUseCase) Execute(commissionID int64, maxPrice float64) error {
	sc, err := uc.repo.FindByID(commissionID)
	if err != nil {
		return err
	}
	if err := sc.UpdateMaxPrice(maxPrice); err != nil {
		return err
	}
	return uc.repo.Save(sc)
}
