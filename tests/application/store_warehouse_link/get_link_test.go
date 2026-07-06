package storewarehouselink_test

import (
	"testing"

	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
	"stock-service/internal/application/store_warehouse_link"
)

type getLinkInMemoryRepo struct {
	links  map[int64]*domainstorewarehouselink.StoreWarehouseLink
	nextID int64
}

func newGetLinkInMemoryRepo() *getLinkInMemoryRepo {
	return &getLinkInMemoryRepo{
		links:  make(map[int64]*domainstorewarehouselink.StoreWarehouseLink),
		nextID: 1,
	}
}

func (r *getLinkInMemoryRepo) Save(swl *domainstorewarehouselink.StoreWarehouseLink) error {
	if swl.ID == 0 {
		swl.ID = r.nextID
		r.nextID++
	}
	r.links[swl.ID] = swl
	return nil
}

func (r *getLinkInMemoryRepo) FindByID(id int64) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl, ok := r.links[id]
	if !ok {
		return nil, nil
	}
	return swl, nil
}

func (r *getLinkInMemoryRepo) FindAll(_ domainstorewarehouselink.WarehouseLinkFilter) ([]*domainstorewarehouselink.StoreWarehouseLink, int, error) {
	var result []*domainstorewarehouselink.StoreWarehouseLink
	for _, swl := range r.links {
		result = append(result, swl)
	}
	return result, len(result), nil
}

func (r *getLinkInMemoryRepo) Delete(id int64) error {
	delete(r.links, id)
	return nil
}

func TestGetLink_Success(t *testing.T) {
	repo := newGetLinkInMemoryRepo()
	uc := storewarehouselink.NewGetLinkUseCase(repo)

	swl := domainstorewarehouselink.NewStoreWarehouseLink(1, 2)
	repo.Save(swl)

	result, err := uc.Execute(storewarehouselink.GetLinkInput{ID: swl.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.StoreID != 1 || result.WarehouseID != 2 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGetLink_NotFound_ReturnsError(t *testing.T) {
	repo := newGetLinkInMemoryRepo()
	uc := storewarehouselink.NewGetLinkUseCase(repo)

	_, err := uc.Execute(storewarehouselink.GetLinkInput{ID: 999})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
