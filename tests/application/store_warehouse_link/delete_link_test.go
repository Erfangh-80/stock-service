package storewarehouselink_test

import (
	"testing"

	domainstorewarehouselink "stock-service/internal/domain/store_warehouse_link"
	"stock-service/internal/application/store_warehouse_link"
)

type deleteLinkInMemoryRepo struct {
	links  map[int64]*domainstorewarehouselink.StoreWarehouseLink
	nextID int64
}

func newDeleteLinkInMemoryRepo() *deleteLinkInMemoryRepo {
	return &deleteLinkInMemoryRepo{
		links:  make(map[int64]*domainstorewarehouselink.StoreWarehouseLink),
		nextID: 1,
	}
}

func (r *deleteLinkInMemoryRepo) Save(swl *domainstorewarehouselink.StoreWarehouseLink) error {
	if swl.ID == 0 {
		swl.ID = r.nextID
		r.nextID++
	}
	r.links[swl.ID] = swl
	return nil
}

func (r *deleteLinkInMemoryRepo) FindByID(id int64) (*domainstorewarehouselink.StoreWarehouseLink, error) {
	swl, ok := r.links[id]
	if !ok {
		return nil, nil
	}
	return swl, nil
}

func (r *deleteLinkInMemoryRepo) FindAll(_ domainstorewarehouselink.WarehouseLinkFilter) ([]*domainstorewarehouselink.StoreWarehouseLink, int, error) {
	var result []*domainstorewarehouselink.StoreWarehouseLink
	for _, swl := range r.links {
		result = append(result, swl)
	}
	return result, len(result), nil
}

func (r *deleteLinkInMemoryRepo) Delete(id int64) error {
	delete(r.links, id)
	return nil
}

func TestDeleteLink_Success(t *testing.T) {
	repo := newDeleteLinkInMemoryRepo()
	uc := storewarehouselink.NewDeleteLinkUseCase(repo)

	swl := domainstorewarehouselink.NewStoreWarehouseLink(1, 2)
	repo.Save(swl)

	err := uc.Execute(storewarehouselink.DeleteLinkInput{ID: swl.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(swl.ID)
	if saved != nil {
		t.Error("expected nil after delete")
	}
}

func TestDeleteLink_NotFound_ReturnsError(t *testing.T) {
	repo := newDeleteLinkInMemoryRepo()
	uc := storewarehouselink.NewDeleteLinkUseCase(repo)

	err := uc.Execute(storewarehouselink.DeleteLinkInput{ID: 999})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
