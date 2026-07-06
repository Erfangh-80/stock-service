package storeallowedcategory_test

import (
	"testing"

	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
	"stock-service/internal/application/store_allowed_category"
)

type deleteCategoryInMemoryRepo struct {
	categories map[int64]*domainstoreallowedcategory.StoreAllowedCategory
	nextID     int64
}

func newDeleteCategoryInMemoryRepo() *deleteCategoryInMemoryRepo {
	return &deleteCategoryInMemoryRepo{
		categories: make(map[int64]*domainstoreallowedcategory.StoreAllowedCategory),
		nextID:     1,
	}
}

func (r *deleteCategoryInMemoryRepo) Save(sac *domainstoreallowedcategory.StoreAllowedCategory) error {
	if sac.ID == 0 {
		sac.ID = r.nextID
		r.nextID++
	}
	r.categories[sac.ID] = sac
	return nil
}

func (r *deleteCategoryInMemoryRepo) FindByID(id int64) (*domainstoreallowedcategory.StoreAllowedCategory, error) {
	sac, ok := r.categories[id]
	if !ok {
		return nil, nil
	}
	return sac, nil
}

func (r *deleteCategoryInMemoryRepo) FindAll(_ domainstoreallowedcategory.StoreCategoryFilter) ([]*domainstoreallowedcategory.StoreAllowedCategory, int, error) {
	var result []*domainstoreallowedcategory.StoreAllowedCategory
	for _, sac := range r.categories {
		result = append(result, sac)
	}
	return result, len(result), nil
}

func (r *deleteCategoryInMemoryRepo) Delete(id int64) error {
	delete(r.categories, id)
	return nil
}

func TestDeleteStoreCategory_Success(t *testing.T) {
	repo := newDeleteCategoryInMemoryRepo()
	uc := storeallowedcategory.NewDeleteStoreCategoryUseCase(repo)

	sac := domainstoreallowedcategory.NewStoreAllowedCategory(1, 2)
	repo.Save(sac)

	err := uc.Execute(storeallowedcategory.DeleteStoreCategoryInput{ID: sac.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(sac.ID)
	if saved != nil {
		t.Error("expected nil after delete")
	}
}

func TestDeleteStoreCategory_NotFound_ReturnsError(t *testing.T) {
	repo := newDeleteCategoryInMemoryRepo()
	uc := storeallowedcategory.NewDeleteStoreCategoryUseCase(repo)

	err := uc.Execute(storeallowedcategory.DeleteStoreCategoryInput{ID: 999})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
