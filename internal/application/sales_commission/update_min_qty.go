package salescommission

import (
	"stock-service/internal/domain/sales_commission"
)

type UpdateMinQtyUseCase struct {
	repo salescommission.Repository
}

func NewUpdateMinQtyUseCase(repo salescommission.Repository) *UpdateMinQtyUseCase {
	return &UpdateMinQtyUseCase{repo: repo}
}

func (uc *UpdateMinQtyUseCase) Execute(commissionID int64, minQty int) error {
	sc, err := uc.repo.FindByID(commissionID)
	if err != nil {
		return err
	}
	if err := sc.UpdateMinQty(minQty); err != nil {
		return err
	}
	return uc.repo.Save(sc)
}
