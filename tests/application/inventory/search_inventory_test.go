package inventory_test

import (
	"strings"
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
	"stock-service/internal/domain/product"
)

type mockProductRepo struct {
	products []*product.Product
}

func (m *mockProductRepo) FindByID(id int32) (*product.Product, error) {
	for _, p := range m.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockProductRepo) FindByTitle(query string) ([]*product.Product, error) {
	var result []*product.Product
	for _, p := range m.products {
		if strings.Contains(p.TitleFa, query) {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockProductRepo) Save(p *product.Product) error { return nil }
func (m *mockProductRepo) FindAll(filter product.ProductFilter) ([]*product.Product, error) { return m.products, nil }
func (m *mockProductRepo) Count(filter product.ProductFilter) (int, error) { return len(m.products), nil }

func TestSearchInventory_ByProductName(t *testing.T) {
	invRepo := newInmemoryRepository()
	prodRepo := &mockProductRepo{
		products: []*product.Product{
			{ID: 42, TitleFa: "محصول چهل و دو"},
			{ID: 100, TitleFa: "محصول صد"},
		},
	}

	inv1, _ := inventory.NewInventory(1, 1, 42, 100)
	inv2, _ := inventory.NewInventory(1, 1, 100, 200)
	inv3, _ := inventory.NewInventory(2, 1, 999, 300)
	invRepo.Save(inv1)
	invRepo.Save(inv2)
	invRepo.Save(inv3)

	uc := appinventory.NewSearchInventoryUseCase(invRepo, prodRepo)
	result, err := uc.Execute(appinventory.SearchInventoryInput{Query: "چهل و دو", Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Errorf("expected 1 inventory item, got %d", result.Total)
	}
	if result.Items[0].ProductID != 42 {
		t.Errorf("expected product ID 42, got %d", result.Items[0].ProductID)
	}
}

func TestSearchInventory_NoResults(t *testing.T) {
	invRepo := newInmemoryRepository()
	prodRepo := &mockProductRepo{
		products: []*product.Product{
			{ID: 1, TitleFa: "test"},
		},
	}

	uc := appinventory.NewSearchInventoryUseCase(invRepo, prodRepo)
	result, err := uc.Execute(appinventory.SearchInventoryInput{Query: "nonexistent", Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 0 {
		t.Errorf("expected 0 results, got %d", result.Total)
	}
}
