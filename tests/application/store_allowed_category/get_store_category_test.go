package storeallowedcategory_test

import (
	"testing"

	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
	"stock-service/internal/application/store_allowed_category"
)

type getCategoryInMemoryRepo struct {
	categories map[int64]*domainstoreallowedcategory.StoreAllowedCategory
	nextID     int64
}

func newGetCategoryInMemoryRepo() *getCategoryInMemoryRepo {
	return &getCategoryInMemoryRepo{
		categories: make(map[int64]*domainstoreallowedcategory.StoreAllowedCategory),
		nextID:     1,
	}
}

func (r *getCategoryInMemoryRepo) Save(sac *domainstoreallowedcategory.StoreAllowedCategory) error {
	if sac.ID == 0 {
		sac.ID = r.nextID
		r.nextID++
	}
	r.categories[sac.ID] = sac
	return nil
}

func (r *getCategoryInMemoryRepo) FindByID(id int64) (*domainstoreallowedcategory.StoreAllowedCategory, error) {
	sac, ok := r.categories[id]
	if !ok {
		return nil, nil
	}
	return sac, nil
}

func (r *getCategoryInMemoryRepo) FindAll(_ domainstoreallowedcategory.StoreCategoryFilter) ([]*domainstoreallowedcategory.StoreAllowedCategory, int, error) {
	var result []*domainstoreallowedcategory.StoreAllowedCategory
	for _, sac := range r.categories {
		result = append(result, sac)
	}
	return result, len(result), nil
}

func (r *getCategoryInMemoryRepo) Delete(id int64) error {
	delete(r.categories, id)
	return nil
}

func TestGetStoreCategory_Success(t *testing.T) {
	repo := newGetCategoryInMemoryRepo()
	uc := storeallowedcategory.NewGetStoreCategoryUseCase(repo)

	sac := domainstoreallowedcategory.NewStoreAllowedCategory(1, 2)
	repo.Save(sac)

	result, err := uc.Execute(storeallowedcategory.GetStoreCategoryInput{ID: sac.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.StoreID != 1 || result.CategoryID != 2 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGetStoreCategory_NotFound_ReturnsError(t *testing.T) {
	repo := newGetCategoryInMemoryRepo()
	uc := storeallowedcategory.NewGetStoreCategoryUseCase(repo)

	_, err := uc.Execute(storeallowedcategory.GetStoreCategoryInput{ID: 999})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
