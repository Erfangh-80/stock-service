package storewarehouselink_test

import (
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
		return nil, nil
	}
	return swl, nil
}

func (r *changeRelationInMemoryRepo) FindAll(_ domainstorewarehouselink.WarehouseLinkFilter) ([]*domainstorewarehouselink.StoreWarehouseLink, int, error) {
	var result []*domainstorewarehouselink.StoreWarehouseLink
	for _, swl := range r.links {
		result = append(result, swl)
	}
	return result, len(result), nil
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

	result, err := uc.Execute(storewarehouselink.ChangeRelationInput{
		LinkID:       swl.ID,
		RelationType: domainstorewarehouselink.RelationTypePrimary,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.RelationType != domainstorewarehouselink.RelationTypePrimary {
		t.Errorf("expected RelationType %q, got %q", domainstorewarehouselink.RelationTypePrimary, result.RelationType)
	}
}

func TestChangeRelation_NotFound_ReturnsError(t *testing.T) {
	repo := newChangeRelationInMemoryRepo()
	uc := storewarehouselink.NewChangeRelationUseCase(repo)

	_, err := uc.Execute(storewarehouselink.ChangeRelationInput{
		LinkID:       999,
		RelationType: domainstorewarehouselink.RelationTypePrimary,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestChangeRelation_InvalidType_ReturnsError(t *testing.T) {
	repo := newChangeRelationInMemoryRepo()
	uc := storewarehouselink.NewChangeRelationUseCase(repo)

	swl := domainstorewarehouselink.NewStoreWarehouseLink(1, 2)
	repo.Save(swl)

	_, err := uc.Execute(storewarehouselink.ChangeRelationInput{
		LinkID:       swl.ID,
		RelationType: "invalid",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
