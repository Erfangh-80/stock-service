package referenceprice_test

import (
	"testing"

	referencepriceapp "stock-service/internal/application/reference_price"
	"stock-service/internal/domain/reference_price"
	inventorydomain "stock-service/internal/domain/inventory"
)

type mockInventoryRepo struct {
	items []*inventorydomain.Inventory
}

func (m *mockInventoryRepo) Save(inv *inventorydomain.Inventory) error { return nil }
func (m *mockInventoryRepo) FindByID(id int64) (*inventorydomain.Inventory, error) { return nil, nil }
func (m *mockInventoryRepo) FindAll() ([]*inventorydomain.Inventory, error) { return m.items, nil }
func (m *mockInventoryRepo) Delete(id int64) error { return nil }

func TestValidateReferencePrice_NoInventory(t *testing.T) {
	repo := populateRepo()
	invRepo := &mockInventoryRepo{items: nil}
	uc := referencepriceapp.NewValidateReferencePriceUseCase(repo, invRepo)

	result, err := uc.Execute(referencepriceapp.ValidateReferencePriceInput{ProductID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Comparison != "no_inventory" {
		t.Errorf("expected 'no_inventory', got %q", result.Comparison)
	}
	if result.ReferencePrice != 99.99 {
		t.Errorf("expected 99.99, got %f", result.ReferencePrice)
	}
}

func TestValidateReferencePrice_NotFound(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	invRepo := &mockInventoryRepo{items: nil}
	uc := referencepriceapp.NewValidateReferencePriceUseCase(repo, invRepo)

	_, err := uc.Execute(referencepriceapp.ValidateReferencePriceInput{ProductID: 999})
	if err != referenceprice.ErrReferencePriceNotFound {
		t.Errorf("expected ErrReferencePriceNotFound, got %v", err)
	}
}

func TestValidateReferencePrice_WithinRange(t *testing.T) {
	repo := populateRepo()
	invRepo := &mockInventoryRepo{
		items: []*inventorydomain.Inventory{
			{ProductID: 1, BasePrice: 100},
		},
	}
	uc := referencepriceapp.NewValidateReferencePriceUseCase(repo, invRepo)

	result, err := uc.Execute(referencepriceapp.ValidateReferencePriceInput{ProductID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Comparison != "within_range" {
		t.Errorf("expected 'within_range', got %q", result.Comparison)
	}
	if len(result.BasePrices) != 1 || result.BasePrices[0] != 100 {
		t.Errorf("unexpected base prices: %v", result.BasePrices)
	}
	if result.InventoryCount != 1 {
		t.Errorf("expected 1 inventory, got %d", result.InventoryCount)
	}
}
