package storewarehouselink_test

import (
	"errors"
	"testing"

	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
	"stock-service/internal/application/store_warehouse_link"
)

type changeRelationInMemoryRepo struct {
	links  map[int64]*domainstorewarehouselink.StoreWarehouseLink
	nextID int64
}

func newChangeRelationInMemoryRepo() *changeRelationInMemoryRepo {
	return &changeRelationInMemoryRepo{
		links:  make(map[int64]*domainstorewarehouselink.StoreWarehouseLink),
		nextID: 1,
	}
}

func (r *changeRelationInMemoryRepo) Save(swl *domainstorewarehouselink.StoreWarehouseLink) error {
	if swl.ID == 0 {
		swl.ID = r.nextID
		r.nextID++
	}
	r.links[swl.ID] = swl
	return nil
}

func (r *changeRelationInMemoryRepo) FindByID(id int64) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl, ok := r.links[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return swl, nil
}

func (r *changeRelationInMemoryRepo) Delete(id int64) error {
	delete(r.links, id)
	return nil
}

func TestChangeRelation_Success(t *testing.T) {
	repo := newChangeRelationInMemoryRepo()
	uc := storewarehouselink.NewChangeRelationUseCase(repo)

	swl := domainstorewarehouselink.NewStoreWarehouseLink(1, 2)
	repo.Save(swl)

	newType := domainstorewarehouselink.RelationType("secondary")
	err := uc.Execute(swl.ID, newType)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(swl.ID)
	if saved.RelationType != newType {
		t.Errorf("expected RelationType %q, got %q", newType, saved.RelationType)
	}
}

func TestChangeRelation_NotFound_ReturnsError(t *testing.T) {
	repo := newChangeRelationInMemoryRepo()
	uc := storewarehouselink.NewChangeRelationUseCase(repo)

	err := uc.Execute(999, domainstorewarehouselink.RelationTypePrimary)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
