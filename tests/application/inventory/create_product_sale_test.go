package inventory_test

import (
	"strings"
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
	"stock-service/internal/domain/product"
)

type inmemoryProductRepo struct {
	products map[int32]*product.Product
}

func newInmemoryProductRepo() *inmemoryProductRepo {
	return &inmemoryProductRepo{products: make(map[int32]*product.Product)}
}

func (r *inmemoryProductRepo) FindByID(id int32) (*product.Product, error) {
	p, ok := r.products[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *inmemoryProductRepo) Save(p *product.Product) error {
	r.products[p.ID] = p
	return nil
}

func (r *inmemoryProductRepo) FindByTitle(query string) ([]*product.Product, error) {
	var result []*product.Product
	for _, p := range r.products {
		if strings.Contains(p.TitleFa, query) {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *inmemoryProductRepo) FindAll(filter product.ProductFilter) ([]*product.Product, error) {
	var result []*product.Product
	for _, p := range r.products {
		result = append(result, p)
	}
	return result, nil
}

func (r *inmemoryProductRepo) Count(filter product.ProductFilter) (int, error) {
	return len(r.products), nil
}

func TestCreateInventoryUseCase_Success(t *testing.T) {
	invRepo := newInmemoryRepository()
	prodRepo := newInmemoryProductRepo()

	p, err := product.NewProduct("test", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	p.ID = 30
	prodRepo.Save(p)

	uc := appinventory.NewCreateInventoryUseCase(invRepo, prodRepo)

	input := appinventory.CreateInventoryInput{
		StoreID:     10,
		WarehouseID: 20,
		ProductID:   30,
		BasePrice:   99.99,
	}

	sale, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale.ID == 0 {
		t.Fatal("expected sale ID to be assigned")
	}
	if sale.StoreID != 10 {
		t.Fatalf("expected StoreID 10, got %d", sale.StoreID)
	}
	if sale.WarehouseID != 20 {
		t.Fatalf("expected WarehouseID 20, got %d", sale.WarehouseID)
	}
	if sale.ProductID != 30 {
		t.Fatalf("expected ProductID 30, got %d", sale.ProductID)
	}
	if sale.BasePrice != 99.99 {
		t.Fatalf("expected BasePrice 99.99, got %f", sale.BasePrice)
	}
}

func TestCreateInventoryUseCase_ProductNotFound(t *testing.T) {
	invRepo := newInmemoryRepository()
	prodRepo := newInmemoryProductRepo()

	uc := appinventory.NewCreateInventoryUseCase(invRepo, prodRepo)

	input := appinventory.CreateInventoryInput{
		StoreID:     10,
		WarehouseID: 20,
		ProductID:   999,
		BasePrice:   99.99,
	}

	_, err := uc.Execute(input)
	if err != product.ErrProductNotFound {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestCreateInventoryUseCase_InvalidBasePrice(t *testing.T) {
	invRepo := newInmemoryRepository()
	prodRepo := newInmemoryProductRepo()

	p, _ := product.NewProduct("test", 1, 1)
	p.ID = 30
	prodRepo.Save(p)

	uc := appinventory.NewCreateInventoryUseCase(invRepo, prodRepo)

	input := appinventory.CreateInventoryInput{
		StoreID:     10,
		WarehouseID: 20,
		ProductID:   30,
		BasePrice:   0,
	}

	_, err := uc.Execute(input)
	if err != inventory.ErrInvalidBasePrice {
		t.Fatalf("expected ErrInvalidBasePrice, got %v", err)
	}
}
