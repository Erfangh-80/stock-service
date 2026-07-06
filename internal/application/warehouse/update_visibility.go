package warehouse

import (
	"stock-service/internal/domain/warehouse"
)

type UpdateVisibilityUseCase struct {
	repo warehouse.Repository
}

func NewUpdateVisibilityUseCase(repo warehouse.Repository) *UpdateVisibilityUseCase {
	return &UpdateVisibilityUseCase{repo: repo}
}

func (uc *UpdateVisibilityUseCase) Execute(warehouseID int64, isPublic bool) error {
	w, err := uc.repo.FindByID(warehouseID)
	if err != nil {
		return err
	}
	if w == nil {
		return warehouse.ErrWarehouseNotFound
	}
	if isPublic {
		w.MakePublic()
	} else {
		w.MakePrivate()
	}
	return uc.repo.Save(w)
}
