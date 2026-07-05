package storeallowedcategory

import (
	"stock-service/internal/domain/store_allowed_category"
)

type CreateCategoryUseCase struct {
	repo storeallowedcategory.Repository
}

func NewCreateCategoryUseCase(repo storeallowedcategory.Repository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{repo: repo}
}

func (uc *CreateCategoryUseCase) Execute(storeID, categoryID int64) (*storeallowedcategory.StoreAllowedCategory, error) {
	sac := storeallowedcategory.NewStoreAllowedCategory(storeID, categoryID)
	if err := uc.repo.Save(sac); err != nil {
		return nil, err
	}
	return sac, nil
}
