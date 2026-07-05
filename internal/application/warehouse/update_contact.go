package warehouse

import (
	"stock-service/internal/domain/warehouse"
)

type UpdateContactUseCase struct {
	repo warehouse.Repository
}

func NewUpdateContactUseCase(repo warehouse.Repository) *UpdateContactUseCase {
	return &UpdateContactUseCase{repo: repo}
}

func (uc *UpdateContactUseCase) Execute(warehouseID int64, phone, contactPhone *string, collectionMethod string) error {
	w, err := uc.repo.FindByID(warehouseID)
	if err != nil {
		return err
	}
	w.UpdatePhone(phone)
	w.UpdateContactPhone(contactPhone)
	w.UpdateCollectionMethod(collectionMethod)
	return uc.repo.Save(w)
}
