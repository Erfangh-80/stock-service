package storeallowedcategory_test

import (
	"testing"

	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
	"stock-service/internal/application/store_allowed_category"
)

type listCategoryInMemoryRepo struct {
	categories map[int64]*domainstoreallowedcategory.StoreAllowedCategory
	nextID     int64
}

func newListCategoryInMemoryRepo() *listCategoryInMemoryRepo {
	return &listCategoryInMemoryRepo{
		categories: make(map[int64]*domainstoreallowedcategory.StoreAllowedCategory),
		nextID:     1,
	}
}

func (r *listCategoryInMemoryRepo) Save(sac *domainstoreallowedcategory.StoreAllowedCategory) error {
	if sac.ID == 0 {
		sac.ID = r.nextID
		r.nextID++
	}
	r.categories[sac.ID] = sac
	return nil
}

func (r *listCategoryInMemoryRepo) FindByID(id int64) (*domainstoreallowedcategory.StoreAllowedCategory, error) {
	sac, ok := r.categories[id]
	if !ok {
		return nil, nil
	}
	return sac, nil
}

func (r *listCategoryInMemoryRepo) FindAll(filter domainstoreallowedcategory.StoreCategoryFilter) ([]*domainstoreallowedcategory.StoreAllowedCategory, int, error) {
	var matched []*domainstoreallowedcategory.StoreAllowedCategory
	for _, sac := range r.categories {
		if filter.StoreID != nil && sac.StoreID != *filter.StoreID {
			continue
		}
		matched = append(matched, sac)
	}
	return matched, len(matched), nil
}

func (r *listCategoryInMemoryRepo) Delete(id int64) error {
	delete(r.categories, id)
	return nil
}

func TestListStoreCategories_Empty(t *testing.T) {
	repo := newListCategoryInMemoryRepo()
	uc := storeallowedcategory.NewListStoreCategoriesUseCase(repo)

	result, err := uc.Execute(storeallowedcategory.ListStoreCategoriesInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Categories) != 0 {
		t.Errorf("expected 0 categories, got %d", len(result.Categories))
	}
}

func TestListStoreCategories_All(t *testing.T) {
	repo := newListCategoryInMemoryRepo()
	uc := storeallowedcategory.NewListStoreCategoriesUseCase(repo)

	repo.Save(domainstoreallowedcategory.NewStoreAllowedCategory(1, 10))
	repo.Save(domainstoreallowedcategory.NewStoreAllowedCategory(1, 20))
	repo.Save(domainstoreallowedcategory.NewStoreAllowedCategory(2, 30))

	result, err := uc.Execute(storeallowedcategory.ListStoreCategoriesInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Categories) != 3 {
		t.Errorf("expected 3 categories, got %d", len(result.Categories))
	}
}

func TestListStoreCategories_FilterByStore(t *testing.T) {
	repo := newListCategoryInMemoryRepo()
	uc := storeallowedcategory.NewListStoreCategoriesUseCase(repo)

	repo.Save(domainstoreallowedcategory.NewStoreAllowedCategory(1, 10))
	repo.Save(domainstoreallowedcategory.NewStoreAllowedCategory(1, 20))
	repo.Save(domainstoreallowedcategory.NewStoreAllowedCategory(2, 30))

	storeID := int64(1)
	result, err := uc.Execute(storeallowedcategory.ListStoreCategoriesInput{
		StoreID: &storeID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(result.Categories))
	}
}
