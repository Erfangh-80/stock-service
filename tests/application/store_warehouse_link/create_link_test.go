package storewarehouselink_test

import (
	"errors"
	"testing"

	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
	"stock-service/internal/application/store_warehouse_link"
)

type createLinkInMemoryRepo struct {
	links  map[int64]*domainstorewarehouselink.StoreWarehouseLink
	nextID int64
}

func newCreateLinkInMemoryRepo() *createLinkInMemoryRepo {
	return &createLinkInMemoryRepo{
		links:  make(map[int64]*domainstorewarehouselink.StoreWarehouseLink),
		nextID: 1,
	}
}

func (r *createLinkInMemoryRepo) Save(swl *domainstorewarehouselink.StoreWarehouseLink) error {
	if swl.ID == 0 {
		swl.ID = r.nextID
		r.nextID++
	}
	r.links[swl.ID] = swl
	return nil
}

func (r *createLinkInMemoryRepo) FindByID(id int64) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl, ok := r.links[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return swl, nil
}

func (r *createLinkInMemoryRepo) Delete(id int64) error {
	delete(r.links, id)
	return nil
}

func TestCreateLink_Success(t *testing.T) {
	repo := newCreateLinkInMemoryRepo()
	uc := storewarehouselink.NewCreateLinkUseCase(repo)

	swl, err := uc.Execute(10, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if swl.ID == 0 {
		t.Error("expected ID to be set")
	}
	if swl.StoreID != 10 {
		t.Errorf("expected StoreID %d, got %d", 10, swl.StoreID)
	}
	if swl.WarehouseID != 20 {
		t.Errorf("expected WarehouseID %d, got %d", 20, swl.WarehouseID)
	}
	if swl.RelationType != domainstorewarehouselink.RelationTypePrimary {
		t.Errorf("expected RelationType %q, got %q", domainstorewarehouselink.RelationTypePrimary, swl.RelationType)
	}
}
