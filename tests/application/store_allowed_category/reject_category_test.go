package storeallowedcategory_test

import (
	"errors"
	"testing"

	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
	"stock-service/internal/application/store_allowed_category"
)

type rejectCategoryInMemoryRepo struct {
	categories map[int64]*domainstoreallowedcategory.StoreAllowedCategory
	nextID     int64
}

func newRejectCategoryInMemoryRepo() *rejectCategoryInMemoryRepo {
	return &rejectCategoryInMemoryRepo{
		categories: make(map[int64]*domainstoreallowedcategory.StoreAllowedCategory),
		nextID:     1,
	}
}

func (r *rejectCategoryInMemoryRepo) Save(sac *domainstoreallowedcategory.StoreAllowedCategory) error {
	if sac.ID == 0 {
		sac.ID = r.nextID
		r.nextID++
	}
	r.categories[sac.ID] = sac
	return nil
}

func (r *rejectCategoryInMemoryRepo) FindByID(id int64) (*domainstoreallowedcategory.StoreAllowedCategory, error) {
	sac, ok := r.categories[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return sac, nil
}

func (r *rejectCategoryInMemoryRepo) Delete(id int64) error {
	delete(r.categories, id)
	return nil
}

func TestRejectCategory_Success(t *testing.T) {
	repo := newRejectCategoryInMemoryRepo()
	uc := storeallowedcategory.NewRejectCategoryUseCase(repo)

	sac := domainstoreallowedcategory.NewStoreAllowedCategory(1, 2)
	repo.Save(sac)

	err := uc.Execute(sac.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(sac.ID)
	if saved.Status != domainstoreallowedcategory.StatusRejected {
		t.Errorf("expected Status %q, got %q", domainstoreallowedcategory.StatusRejected, saved.Status)
	}
}

func TestRejectCategory_NotFound_ReturnsError(t *testing.T) {
	repo := newRejectCategoryInMemoryRepo()
	uc := storeallowedcategory.NewRejectCategoryUseCase(repo)

	err := uc.Execute(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
