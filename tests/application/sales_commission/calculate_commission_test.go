package salescommission_test

import (
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	inventorydomain "stock-service/internal/domain/inventory"
	"stock-service/internal/application/sales_commission"
)

type mockInventoryRepo struct {
	items map[int64]*inventorydomain.Inventory
}

func newMockInventoryRepo() *mockInventoryRepo {
	return &mockInventoryRepo{items: make(map[int64]*inventorydomain.Inventory)}
}

func (m *mockInventoryRepo) Save(inv *inventorydomain.Inventory) error {
	if inv.ID == 0 {
		inv.ID = int64(len(m.items) + 1)
	}
	m.items[inv.ID] = inv
	return nil
}
func (m *mockInventoryRepo) FindByID(id int64) (*inventorydomain.Inventory, error) { return m.items[id], nil }
func (m *mockInventoryRepo) FindAll() ([]*inventorydomain.Inventory, error) {
	var result []*inventorydomain.Inventory
	for _, v := range m.items {
		result = append(result, v)
	}
	return result, nil
}
func (m *mockInventoryRepo) Delete(id int64) error { delete(m.items, id); return nil }

func TestCalculateCommission_BasePrice(t *testing.T) {
	commRepo := newInMemorySalesCommissionRepo()
	invRepo := newMockInventoryRepo()

	inv, _ := inventorydomain.NewInventory(1, 1, 1, 500)
	invRepo.Save(inv)

	sc, _ := domainsalescommission.NewSalesCommission(inv.ID, 1, domainsalescommission.SaleModelRetail, 10, 100)
	commRepo.Save(sc)

	uc := salescommission.NewCalculateCommissionUseCase(commRepo, invRepo)
	result, err := uc.Execute(salescommission.CalculateCommissionInput{InventoryID: inv.ID, Quantity: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.PriceSource != "base_price" {
		t.Errorf("expected base_price, got %s", result.PriceSource)
	}
	if result.BasePriceUsed != 500 {
		t.Errorf("expected 500, got %f", result.BasePriceUsed)
	}
	if result.CommissionAmt != 100 { // 500 * 10% * 2 = 100
		t.Errorf("expected 100, got %f", result.CommissionAmt)
	}
}

func TestCalculateCommission_FinalPrice(t *testing.T) {
	commRepo := newInMemorySalesCommissionRepo()
	invRepo := newMockInventoryRepo()

	fp := 400.0
	pid := int64(1)
	inv, _ := inventorydomain.NewInventory(1, 1, 1, 500)
	inv.FinalPrice = &fp
	inv.PromotionID = &pid
	invRepo.Save(inv)

	sc, _ := domainsalescommission.NewSalesCommission(inv.ID, 1, domainsalescommission.SaleModelRetail, 10, 100)
	commRepo.Save(sc)

	uc := salescommission.NewCalculateCommissionUseCase(commRepo, invRepo)
	result, err := uc.Execute(salescommission.CalculateCommissionInput{InventoryID: inv.ID, Quantity: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.PriceSource != "final_price" {
		t.Errorf("expected final_price, got %s", result.PriceSource)
	}
	if result.BasePriceUsed != 400 {
		t.Errorf("expected 400, got %f", result.BasePriceUsed)
	}
	if result.CommissionAmt != 40 { // 400 * 10% * 1 = 40
		t.Errorf("expected 40, got %f", result.CommissionAmt)
	}
}

func TestCalculateCommission_CommissionNotFound(t *testing.T) {
	commRepo := newInMemorySalesCommissionRepo()
	invRepo := newMockInventoryRepo()

	uc := salescommission.NewCalculateCommissionUseCase(commRepo, invRepo)
	_, err := uc.Execute(salescommission.CalculateCommissionInput{InventoryID: 999, Quantity: 1})
	if err != domainsalescommission.ErrCommissionNotFound {
		t.Errorf("expected ErrCommissionNotFound, got %v", err)
	}
}

func TestCalculateCommission_InventoryNotFound(t *testing.T) {
	commRepo := newInMemorySalesCommissionRepo()
	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 10, 100)
	sc.InventoryID = 999
	commRepo.Save(sc)

	invRepo := newMockInventoryRepo()
	uc := salescommission.NewCalculateCommissionUseCase(commRepo, invRepo)
	_, err := uc.Execute(salescommission.CalculateCommissionInput{InventoryID: 999, Quantity: 1})
	if err != inventorydomain.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}
