package warehouse_test

import (
	"errors"
	"testing"

	domainwarehouse "stock-service/internal/domain/warehouse"
	"stock-service/internal/application/warehouse"
)

type updateContactInMemoryRepo struct {
	warehouses map[int64]*domainwarehouse.Warehouse
	nextID     int64
}

func newUpdateContactInMemoryRepo() *updateContactInMemoryRepo {
	return &updateContactInMemoryRepo{
		warehouses: make(map[int64]*domainwarehouse.Warehouse),
		nextID:     1,
	}
}

func (r *updateContactInMemoryRepo) Save(w *domainwarehouse.Warehouse) error {
	if w.ID == 0 {
		w.ID = r.nextID
		r.nextID++
	}
	r.warehouses[w.ID] = w
	return nil
}

func (r *updateContactInMemoryRepo) FindByID(id int64) (*domainwarehouse.Warehouse, error) {
	w, ok := r.warehouses[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return w, nil
}

func (r *updateContactInMemoryRepo) FindAll(filter domainwarehouse.WarehouseFilter) ([]*domainwarehouse.Warehouse, int, error) {
	var result []*domainwarehouse.Warehouse
	for _, w := range r.warehouses {
		result = append(result, w)
	}
	return result, len(result), nil
}

func (r *updateContactInMemoryRepo) Delete(id int64) error {
	delete(r.warehouses, id)
	return nil
}

func TestUpdateContact_Success(t *testing.T) {
	repo := newUpdateContactInMemoryRepo()
	uc := warehouse.NewUpdateContactUseCase(repo)

	w, _ := domainwarehouse.NewWarehouse(1, "Test Warehouse")
	repo.Save(w)

	phone := "123-456-7890"
	contactPhone := "098-765-4321"
	collectionMethod := "pickup"

	err := uc.Execute(w.ID, &phone, &contactPhone, collectionMethod)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(w.ID)
	if saved.Phone == nil || *saved.Phone != phone {
		t.Errorf("expected Phone %q, got %v", phone, saved.Phone)
	}
	if saved.ContactPhone == nil || *saved.ContactPhone != contactPhone {
		t.Errorf("expected ContactPhone %q, got %v", contactPhone, saved.ContactPhone)
	}
	if saved.CollectionMethod != collectionMethod {
		t.Errorf("expected CollectionMethod %q, got %q", collectionMethod, saved.CollectionMethod)
	}
}

func TestUpdateContact_NotFound_ReturnsError(t *testing.T) {
	repo := newUpdateContactInMemoryRepo()
	uc := warehouse.NewUpdateContactUseCase(repo)

	phone := "123"
	contactPhone := "456"

	err := uc.Execute(999, &phone, &contactPhone, "delivery")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
