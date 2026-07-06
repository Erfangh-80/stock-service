package warehouse_test

import (
	"errors"
	"testing"

	domainwarehouse "stock-service/internal/domain/warehouse"
	"stock-service/internal/application/warehouse"
)

type createWarehouseInMemoryRepo struct {
	warehouses map[int64]*domainwarehouse.Warehouse
	nextID     int64
}

func newCreateWarehouseInMemoryRepo() *createWarehouseInMemoryRepo {
	return &createWarehouseInMemoryRepo{
		warehouses: make(map[int64]*domainwarehouse.Warehouse),
		nextID:     1,
	}
}

func (r *createWarehouseInMemoryRepo) Save(w *domainwarehouse.Warehouse) error {
	if w.ID == 0 {
		w.ID = r.nextID
		r.nextID++
	}
	r.warehouses[w.ID] = w
	return nil
}

func (r *createWarehouseInMemoryRepo) FindByID(id int64) (*domainwarehouse.Warehouse, error) {
	w, ok := r.warehouses[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return w, nil
}

func (r *createWarehouseInMemoryRepo) FindAll(filter domainwarehouse.WarehouseFilter) ([]*domainwarehouse.Warehouse, int, error) {
	var result []*domainwarehouse.Warehouse
	for _, w := range r.warehouses {
		result = append(result, w)
	}
	return result, len(result), nil
}

func (r *createWarehouseInMemoryRepo) Delete(id int64) error {
	delete(r.warehouses, id)
	return nil
}

func TestCreateWarehouse_Success(t *testing.T) {
	repo := newCreateWarehouseInMemoryRepo()
	uc := warehouse.NewCreateWarehouseUseCase(repo)

	w, err := uc.Execute(42, "Main Warehouse")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.ID == 0 {
		t.Error("expected ID to be set")
	}
	if w.CreatedByUserID != 42 {
		t.Errorf("expected CreatedByUserID %d, got %d", 42, w.CreatedByUserID)
	}
	if w.WarehouseName != "Main Warehouse" {
		t.Errorf("expected WarehouseName %q, got %q", "Main Warehouse", w.WarehouseName)
	}
	if w.IsPublic {
		t.Error("expected IsPublic to be false")
	}
}

func TestCreateWarehouse_EmptyName_ReturnsError(t *testing.T) {
	repo := newCreateWarehouseInMemoryRepo()
	uc := warehouse.NewCreateWarehouseUseCase(repo)

	_, err := uc.Execute(1, "")
	if err != domainwarehouse.ErrWarehouseNameRequired {
		t.Errorf("expected %v, got %v", domainwarehouse.ErrWarehouseNameRequired, err)
	}
}
