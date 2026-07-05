package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestListInventory_NoFilters(t *testing.T) {
	repo := newInmemoryRepository()
	inv1, _ := inventory.NewInventory(1, 1, 42, 100)
	inv2, _ := inventory.NewInventory(2, 1, 42, 200)
	repo.Save(inv1)
	repo.Save(inv2)

	uc := appinventory.NewListInventoryUseCase(repo)
	result, err := uc.Execute(appinventory.ListInventoryInput{Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 {
		t.Errorf("expected 2 items, got %d", result.Total)
	}
}

func TestListInventory_FilterByStore(t *testing.T) {
	repo := newInmemoryRepository()
	inv1, _ := inventory.NewInventory(1, 1, 42, 100)
	inv2, _ := inventory.NewInventory(2, 1, 42, 200)
	repo.Save(inv1)
	repo.Save(inv2)

	storeID := int64(1)
	uc := appinventory.NewListInventoryUseCase(repo)
	result, err := uc.Execute(appinventory.ListInventoryInput{StoreID: &storeID, Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Errorf("expected 1 item, got %d", result.Total)
	}
}

func TestListInventory_Pagination(t *testing.T) {
	repo := newInmemoryRepository()
	for i := 0; i < 5; i++ {
		inv, _ := inventory.NewInventory(int64(i+1), 1, 42, 100)
		repo.Save(inv)
	}

	uc := appinventory.NewListInventoryUseCase(repo)
	result, err := uc.Execute(appinventory.ListInventoryInput{Page: 1, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 5 {
		t.Errorf("expected total 5, got %d", result.Total)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items on page 1, got %d", len(result.Items))
	}

	result2, err := uc.Execute(appinventory.ListInventoryInput{Page: 3, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(result2.Items) != 1 {
		t.Errorf("expected 1 item on page 3, got %d", len(result2.Items))
	}
}

func TestListInventory_PageOutOfRange(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	repo.Save(inv)

	uc := appinventory.NewListInventoryUseCase(repo)
	result, err := uc.Execute(appinventory.ListInventoryInput{Page: 10, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
}
