package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestUpdateInventoryUseCase_Success(t *testing.T) {
	repo := newInmemoryRepository()
	repo.Save(&inventory.Inventory{ID: 1})
	uc := appinventory.NewUpdateInventoryUseCase(repo)

	maxQty := 100
	input := appinventory.UpdateInventoryInput{
		SaleID:       1,
		InstantQty:   50,
		ScheduledQty: map[string]int{"2026-07-10": 30},
		MinOrderQty:  5,
		MaxOrderQty:  &maxQty,
	}

	sale, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale.InstantQty != 50 {
		t.Fatalf("expected InstantQty 50, got %d", sale.InstantQty)
	}
	if sale.ScheduledQty["2026-07-10"] != 30 {
		t.Fatalf("expected ScheduledQty[2026-07-10] 30, got %d", sale.ScheduledQty["2026-07-10"])
	}
	if sale.MinOrderQty != 5 {
		t.Fatalf("expected MinOrderQty 5, got %d", sale.MinOrderQty)
	}
	if sale.MaxOrderQty == nil || *sale.MaxOrderQty != 100 {
		t.Fatalf("expected MaxOrderQty 100, got %v", sale.MaxOrderQty)
	}
}

func TestUpdateInventoryUseCase_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	uc := appinventory.NewUpdateInventoryUseCase(repo)

	input := appinventory.UpdateInventoryInput{
		SaleID:      999,
		InstantQty:  10,
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != inventory.ErrInventoryNotFound {
		t.Fatalf("expected errNotFound, got %v", err)
	}
}

func TestUpdateInventoryUseCase_InvalidInstantQty(t *testing.T) {
	repo := newInmemoryRepository()
	repo.Save(&inventory.Inventory{ID: 1})
	uc := appinventory.NewUpdateInventoryUseCase(repo)

	input := appinventory.UpdateInventoryInput{
		SaleID:      1,
		InstantQty:  -1,
		MinOrderQty: 1,
	}

	_, err := uc.Execute(input)
	if err != inventory.ErrInvalidQuantity {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}
}
