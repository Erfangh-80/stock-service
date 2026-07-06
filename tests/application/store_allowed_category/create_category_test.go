package storeallowedcategory_test

import (
	"testing"

	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
	"stock-service/internal/application/store_allowed_category"
)

type createCategoryInMemoryRepo struct {
	categories map[int64]*domainstoreallowedcategory.StoreAllowedCategory
	nextID     int64
}

func newCreateCategoryInMemoryRepo() *createCategoryInMemoryRepo {
	return &createCategoryInMemoryRepo{
		categories: make(map[int64]*domainstoreallowedcategory.StoreAllowedCategory),
		nextID:     1,
	}
}

func (r *createCategoryInMemoryRepo) Save(sac *domainstoreallowedcategory.StoreAllowedCategory) error {
	if sac.ID == 0 {
		sac.ID = r.nextID
		r.nextID++
	}
	r.categories[sac.ID] = sac
	return nil
}

func (r *createCategoryInMemoryRepo) FindByID(id int64) (*domainstoreallowedcategory.StoreAllowedCategory, error) {
	sac, ok := r.categories[id]
	if !ok {
		return nil, nil
	}
	return sac, nil
}

func (r *createCategoryInMemoryRepo) FindAll(_ domainstoreallowedcategory.StoreCategoryFilter) ([]*domainstoreallowedcategory.StoreAllowedCategory, int, error) {
	var result []*domainstoreallowedcategory.StoreAllowedCategory
	for _, sac := range r.categories {
		result = append(result, sac)
	}
	return result, len(result), nil
}

func (r *createCategoryInMemoryRepo) Delete(id int64) error {
	delete(r.categories, id)
	return nil
}

func TestCreateCategory_Success(t *testing.T) {
	repo := newCreateCategoryInMemoryRepo()
	uc := storeallowedcategory.NewCreateCategoryUseCase(repo)

	sac, err := uc.Execute(100, 200)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sac.ID == 0 {
		t.Error("expected ID to be set")
	}
	if sac.StoreID != 100 {
		t.Errorf("expected StoreID %d, got %d", 100, sac.StoreID)
	}
	if sac.CategoryID != 200 {
		t.Errorf("expected CategoryID %d, got %d", 200, sac.CategoryID)
	}
	if sac.Status != domainstoreallowedcategory.StatusPending {
		t.Errorf("expected Status %q, got %q", domainstoreallowedcategory.StatusPending, sac.Status)
	}
}
