package storewarehouselink_test

import (
	"testing"

	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
	"stock-service/internal/application/store_warehouse_link"
)

type listLinksInMemoryRepo struct {
	links  map[int64]*domainstorewarehouselink.StoreWarehouseLink
	nextID int64
}

func newListLinksInMemoryRepo() *listLinksInMemoryRepo {
	return &listLinksInMemoryRepo{
		links:  make(map[int64]*domainstorewarehouselink.StoreWarehouseLink),
		nextID: 1,
	}
}

func (r *listLinksInMemoryRepo) Save(swl *domainstorewarehouselink.StoreWarehouseLink) error {
	if swl.ID == 0 {
		swl.ID = r.nextID
		r.nextID++
	}
	r.links[swl.ID] = swl
	return nil
}

func (r *listLinksInMemoryRepo) FindByID(id int64) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl, ok := r.links[id]
	if !ok {
		return nil, nil
	}
	return swl, nil
}

func (r *listLinksInMemoryRepo) FindAll(filter domainstorewarehouselink.WarehouseLinkFilter) ([]*domainstorewarehouselink.StoreWarehouseLink, int, error) {
	var matched []*domainstorewarehouselink.StoreWarehouseLink
	for _, swl := range r.links {
		if filter.StoreID != nil && swl.StoreID != *filter.StoreID {
			continue
		}
		if filter.WarehouseID != nil && swl.WarehouseID != *filter.WarehouseID {
			continue
		}
		matched = append(matched, swl)
	}
	return matched, len(matched), nil
}

func (r *listLinksInMemoryRepo) Delete(id int64) error {
	delete(r.links, id)
	return nil
}

func TestListLinks_Empty(t *testing.T) {
	repo := newListLinksInMemoryRepo()
	uc := storewarehouselink.NewListLinksUseCase(repo)

	result, err := uc.Execute(storewarehouselink.ListLinksInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Links) != 0 {
		t.Errorf("expected 0 links, got %d", len(result.Links))
	}
}

func TestListLinks_All(t *testing.T) {
	repo := newListLinksInMemoryRepo()
	uc := storewarehouselink.NewListLinksUseCase(repo)

	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(1, 10))
	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(1, 20))
	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(2, 30))

	result, err := uc.Execute(storewarehouselink.ListLinksInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Links) != 3 {
		t.Errorf("expected 3 links, got %d", len(result.Links))
	}
}

func TestListLinks_FilterByStore(t *testing.T) {
	repo := newListLinksInMemoryRepo()
	uc := storewarehouselink.NewListLinksUseCase(repo)

	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(1, 10))
	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(1, 20))
	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(2, 30))

	storeID := int64(1)
	result, err := uc.Execute(storewarehouselink.ListLinksInput{StoreID: &storeID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Links) != 2 {
		t.Errorf("expected 2 links, got %d", len(result.Links))
	}
}

func TestListLinks_FilterByWarehouse(t *testing.T) {
	repo := newListLinksInMemoryRepo()
	uc := storewarehouselink.NewListLinksUseCase(repo)

	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(1, 10))
	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(2, 10))
	repo.Save(domainstorewarehouselink.NewStoreWarehouseLink(1, 20))

	warehouseID := int64(10)
	result, err := uc.Execute(storewarehouselink.ListLinksInput{WarehouseID: &warehouseID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Links) != 2 {
		t.Errorf("expected 2 links, got %d", len(result.Links))
	}
}
