package warehouse_test

import (
	"errors"
	"testing"

	domainwarehouse "stock-service/internal/domain/warehouse"
	"stock-service/internal/application/warehouse"
)

type updateVisibilityInMemoryRepo struct {
	warehouses map[int64]*domainwarehouse.Warehouse
	nextID     int64
}

func newUpdateVisibilityInMemoryRepo() *updateVisibilityInMemoryRepo {
	return &updateVisibilityInMemoryRepo{
		warehouses: make(map[int64]*domainwarehouse.Warehouse),
		nextID:     1,
	}
}

func (r *updateVisibilityInMemoryRepo) Save(w *domainwarehouse.Warehouse) error {
	if w.ID == 0 {
		w.ID = r.nextID
		r.nextID++
	}
	r.warehouses[w.ID] = w
	return nil
}

func (r *updateVisibilityInMemoryRepo) FindByID(id int64) (*domainwarehouse.Warehouse, error) {
	w, ok := r.warehouses[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return w, nil
}

func (r *updateVisibilityInMemoryRepo) FindAll(filter domainwarehouse.WarehouseFilter) ([]*domainwarehouse.Warehouse, int, error) {
	var result []*domainwarehouse.Warehouse
	for _, w := range r.warehouses {
		result = append(result, w)
	}
	return result, len(result), nil
}

func (r *updateVisibilityInMemoryRepo) Delete(id int64) error {
	delete(r.warehouses, id)
	return nil
}

func TestUpdateVisibility_MakePublic(t *testing.T) {
	repo := newUpdateVisibilityInMemoryRepo()
	uc := warehouse.NewUpdateVisibilityUseCase(repo)

	w, _ := domainwarehouse.NewWarehouse(1, "Test Warehouse")
	repo.Save(w)

	err := uc.Execute(w.ID, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(w.ID)
	if !saved.IsPublic {
		t.Error("expected IsPublic to be true")
	}
}

func TestUpdateVisibility_MakePrivate(t *testing.T) {
	repo := newUpdateVisibilityInMemoryRepo()
	uc := warehouse.NewUpdateVisibilityUseCase(repo)

	w, _ := domainwarehouse.NewWarehouse(1, "Test Warehouse")
	w.MakePublic()
	repo.Save(w)

	err := uc.Execute(w.ID, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(w.ID)
	if saved.IsPublic {
		t.Error("expected IsPublic to be false")
	}
}

func TestUpdateVisibility_NotFound_ReturnsError(t *testing.T) {
	repo := newUpdateVisibilityInMemoryRepo()
	uc := warehouse.NewUpdateVisibilityUseCase(repo)

	err := uc.Execute(999, true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
