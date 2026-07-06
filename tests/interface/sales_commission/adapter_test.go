package salescommissioninterface_test

import (
	"testing"
	"time"

	domain "stock-service/internal/domain/sales_commission"
	iface "stock-service/internal/interface"
	salescommissioninterface "stock-service/internal/interface/sales_commission"
)

type mockCreateSalesCommission struct {
	fn func(int64, int64, domain.SaleModel, float64, float64) (*domain.SalesCommission, error)
}

func (m *mockCreateSalesCommission) Execute(inventoryID, categoryCommissionRuleID int64, saleModel domain.SaleModel, ratePercent, minPrice float64) (*domain.SalesCommission, error) {
	return m.fn(inventoryID, categoryCommissionRuleID, saleModel, ratePercent, minPrice)
}

type mockUpdateMaxPrice struct {
	fn func(int64, float64) error
}

func (m *mockUpdateMaxPrice) Execute(commissionID int64, maxPrice float64) error {
	return m.fn(commissionID, maxPrice)
}

type mockUpdateMinQty struct {
	fn func(int64, int) error
}

func (m *mockUpdateMinQty) Execute(commissionID int64, minQty int) error {
	return m.fn(commissionID, minQty)
}

func baseSalesCommission() *domain.SalesCommission {
	return &domain.SalesCommission{
		ID: 1, InventoryID: 100, CategoryCommissionRuleID: 200,
		SaleModel: domain.SaleModelRetail, RatePercent: 10.5, MinPrice: 50,
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestAdapter_Create_Success(t *testing.T) {
	sc := baseSalesCommission()
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{func(inventoryID, categoryCommissionRuleID int64, saleModel domain.SaleModel, ratePercent, minPrice float64) (*domain.SalesCommission, error) {
			return sc, nil
		}},
		&mockUpdateMaxPrice{nil},
		&mockUpdateMinQty{nil},
	)
	resp, err := adapter.Create(salescommissioninterface.CreateSalesCommissionParams{
		InventoryID: 100, CategoryCommissionRuleID: 200,
		SaleModel: "retail", RatePercent: 10.5, MinPrice: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.InventoryID != 100 || resp.CategoryCommissionRuleID != 200 {
		t.Error("unexpected response")
	}
	if resp.SaleModel != "retail" || resp.RatePercent != 10.5 || resp.MinPrice != 50 {
		t.Error("unexpected response values")
	}
	if resp.MinQty != nil {
		t.Error("expected MinQty to be nil")
	}
	if resp.MaxPrice != nil {
		t.Error("expected MaxPrice to be nil")
	}
}

func TestAdapter_Create_ErrInvalidRatePercent(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{func(inventoryID, categoryCommissionRuleID int64, saleModel domain.SaleModel, ratePercent, minPrice float64) (*domain.SalesCommission, error) {
			return nil, domain.ErrInvalidRatePercent
		}},
		&mockUpdateMaxPrice{nil},
		&mockUpdateMinQty{nil},
	)
	_, err := adapter.Create(salescommissioninterface.CreateSalesCommissionParams{
		RatePercent: 150,
	})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAdapter_Create_ErrInvalidMinPrice(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{func(inventoryID, categoryCommissionRuleID int64, saleModel domain.SaleModel, ratePercent, minPrice float64) (*domain.SalesCommission, error) {
			return nil, domain.ErrInvalidMinPrice
		}},
		&mockUpdateMaxPrice{nil},
		&mockUpdateMinQty{nil},
	)
	_, err := adapter.Create(salescommissioninterface.CreateSalesCommissionParams{
		MinPrice: 0,
	})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAdapter_UpdateMaxPrice_Success(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{nil},
		&mockUpdateMaxPrice{func(commissionID int64, maxPrice float64) error {
			if commissionID != 1 || maxPrice != 200 {
				t.Error("unexpected input")
			}
			return nil
		}},
		&mockUpdateMinQty{nil},
	)
	err := adapter.UpdateMaxPrice(salescommissioninterface.UpdateMaxPriceParams{
		CommissionID: 1, MaxPrice: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_UpdateMaxPrice_ErrInvalidMaxPrice(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{nil},
		&mockUpdateMaxPrice{func(commissionID int64, maxPrice float64) error {
			return domain.ErrInvalidMaxPrice
		}},
		&mockUpdateMinQty{nil},
	)
	err := adapter.UpdateMaxPrice(salescommissioninterface.UpdateMaxPriceParams{
		CommissionID: 1, MaxPrice: 10,
	})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAdapter_UpdateMinQty_Success(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{nil},
		&mockUpdateMaxPrice{nil},
		&mockUpdateMinQty{func(commissionID int64, minQty int) error {
			if commissionID != 1 || minQty != 5 {
				t.Error("unexpected input")
			}
			return nil
		}},
	)
	err := adapter.UpdateMinQty(salescommissioninterface.UpdateMinQtyParams{
		CommissionID: 1, MinQty: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_UpdateMinQty_ErrInvalidMinQty(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{nil},
		&mockUpdateMaxPrice{nil},
		&mockUpdateMinQty{func(commissionID int64, minQty int) error {
			return domain.ErrInvalidMinQty
		}},
	)
	err := adapter.UpdateMinQty(salescommissioninterface.UpdateMinQtyParams{
		CommissionID: 1, MinQty: -1,
	})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestAdapter_Create_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{func(inventoryID, categoryCommissionRuleID int64, saleModel domain.SaleModel, ratePercent, minPrice float64) (*domain.SalesCommission, error) {
			return nil, unknownErr{}
		}},
		&mockUpdateMaxPrice{nil},
		&mockUpdateMinQty{nil},
	)
	_, err := adapter.Create(salescommissioninterface.CreateSalesCommissionParams{})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestAdapter_UpdateMaxPrice_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{nil},
		&mockUpdateMaxPrice{func(commissionID int64, maxPrice float64) error {
			return unknownErr{}
		}},
		&mockUpdateMinQty{nil},
	)
	err := adapter.UpdateMaxPrice(salescommissioninterface.UpdateMaxPriceParams{
		CommissionID: 1, MaxPrice: 200,
	})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestAdapter_UpdateMinQty_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{nil},
		&mockUpdateMaxPrice{nil},
		&mockUpdateMinQty{func(commissionID int64, minQty int) error {
			return unknownErr{}
		}},
	)
	err := adapter.UpdateMinQty(salescommissioninterface.UpdateMinQtyParams{
		CommissionID: 1, MinQty: 5,
	})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
