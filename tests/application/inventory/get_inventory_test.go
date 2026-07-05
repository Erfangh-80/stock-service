package inventory_test

import (
	"testing"

	"stock-service/internal/domain/inventory"
	"stock-service/internal/domain/product"
	appinventory "stock-service/internal/application/inventory"
)

func TestGetInventoryUseCase_Success(t *testing.T) {
	repo := newInmemoryRepository()
	prodRepo := newInmemoryProductRepo()

	p, _ := product.NewProduct("test", 1, 1)
	p.ID = 42
	prodRepo.Save(p)

	createUC := appinventory.NewCreateInventoryUseCase(repo, prodRepo)
	getUC := appinventory.NewGetInventoryUseCase(repo)

	created, err := createUC.Execute(appinventory.CreateInventoryInput{
		StoreID: 1, WarehouseID: 1, ProductID: 42, BasePrice: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := getUC.Execute(appinventory.GetInventoryInput{ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
}

func TestGetInventoryUseCase_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	getUC := appinventory.NewGetInventoryUseCase(repo)

	_, err := getUC.Execute(appinventory.GetInventoryInput{ID: 999})
	if err != inventory.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}
