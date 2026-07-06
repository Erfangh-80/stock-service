package salescommissioninterface_test

import (
	"testing"
	"time"

	app "stock-service/internal/application/sales_commission"
	domain "stock-service/internal/domain/sales_commission"
	iface "stock-service/internal/interface"
	salescommissioninterface "stock-service/internal/interface/sales_commission"
)

type mockCreate struct {
	fn func(int64, int64, domain.SaleModel, float64, float64) (*domain.SalesCommission, error)
}
func (m *mockCreate) Execute(a, b int64, c domain.SaleModel, d, e float64) (*domain.SalesCommission, error) { return m.fn(a, b, c, d, e) }

type mockUpdateMaxPrice struct {
	fn func(int64, float64) error
}
func (m *mockUpdateMaxPrice) Execute(a int64, b float64) error { return m.fn(a, b) }

type mockUpdateMinQty struct {
	fn func(int64, int) error
}
func (m *mockUpdateMinQty) Execute(a int64, b int) error { return m.fn(a, b) }

type mockGet struct {
	fn func(app.GetSalesCommissionInput) (*domain.SalesCommission, error)
}
func (m *mockGet) Execute(a app.GetSalesCommissionInput) (*domain.SalesCommission, error) { return m.fn(a) }

type mockGetByInv struct {
	fn func(app.GetByInventorySalesCommissionInput) (*domain.SalesCommission, error)
}
func (m *mockGetByInv) Execute(a app.GetByInventorySalesCommissionInput) (*domain.SalesCommission, error) { return m.fn(a) }

type mockList struct {
	fn func(app.ListSalesCommissionsInput) (*app.ListSalesCommissionsOutput, error)
}
func (m *mockList) Execute(a app.ListSalesCommissionsInput) (*app.ListSalesCommissionsOutput, error) { return m.fn(a) }

type mockDelete struct {
	fn func(app.DeleteSalesCommissionInput) error
}
func (m *mockDelete) Execute(a app.DeleteSalesCommissionInput) error { return m.fn(a) }

type mockCalc struct {
	fn func(app.CalculateCommissionInput) (*app.CommissionCalculation, error)
}
func (m *mockCalc) Execute(a app.CalculateCommissionInput) (*app.CommissionCalculation, error) { return m.fn(a) }

func newTestAdapter(c *mockCreate, u *mockUpdateMaxPrice, q *mockUpdateMinQty,
	g *mockGet, b *mockGetByInv, l *mockList, d *mockDelete, c2 *mockCalc) *salescommissioninterface.Adapter {
	return salescommissioninterface.NewAdapter(c, u, q, g, b, l, d, c2)
}

func baseCommission() *domain.SalesCommission {
	return &domain.SalesCommission{
		ID: 1, InventoryID: 100, CategoryCommissionRuleID: 200,
		SaleModel: domain.SaleModelRetail, RatePercent: 10.5, MinPrice: 50,
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestAdapter_Create_Success(t *testing.T) {
	sc := baseCommission()
	adapter := newTestAdapter(
		&mockCreate{func(int64, int64, domain.SaleModel, float64, float64) (*domain.SalesCommission, error) { return sc, nil }},
		&mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	resp, err := adapter.Create(salescommissioninterface.CreateSalesCommissionParams{
		InventoryID: 100, CategoryCommissionRuleID: 200, SaleModel: "retail", RatePercent: 10.5, MinPrice: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.InventoryID != 100 || resp.RatePercent != 10.5 {
		t.Error("unexpected response")
	}
}

func TestAdapter_Create_ErrInvalidRatePercent(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{func(int64, int64, domain.SaleModel, float64, float64) (*domain.SalesCommission, error) {
			return nil, domain.ErrInvalidRatePercent
		}},
		&mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	_, err := adapter.Create(salescommissioninterface.CreateSalesCommissionParams{RatePercent: 150})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAdapter_Get_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{},
		&mockGet{func(app.GetSalesCommissionInput) (*domain.SalesCommission, error) { return baseCommission(), nil }},
		&mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	resp, err := adapter.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("unexpected response")
	}
}

func TestAdapter_Get_NotFound(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{},
		&mockGet{func(app.GetSalesCommissionInput) (*domain.SalesCommission, error) { return nil, domain.ErrCommissionNotFound }},
		&mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	_, err := adapter.Get(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_GetByInventory_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{},
		&mockGetByInv{func(app.GetByInventorySalesCommissionInput) (*domain.SalesCommission, error) { return baseCommission(), nil }},
		&mockList{}, &mockDelete{}, &mockCalc{},
	)
	resp, err := adapter.GetByInventory(100)
	if err != nil {
		t.Fatal(err)
	}
	if resp.InventoryID != 100 {
		t.Error("unexpected response")
	}
}

func TestAdapter_List_Success(t *testing.T) {
	output := &app.ListSalesCommissionsOutput{
		Commissions: []*domain.SalesCommission{baseCommission()},
		Total: 1, Page: 1, Limit: 20,
	}
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{},
		&mockList{func(app.ListSalesCommissionsInput) (*app.ListSalesCommissionsOutput, error) { return output, nil }},
		&mockDelete{}, &mockCalc{},
	)
	resp, err := adapter.List(salescommissioninterface.ListSalesCommissionsParams{Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Commissions) != 1 || resp.Total != 1 {
		t.Error("unexpected list response")
	}
}

func TestAdapter_Delete_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{},
		&mockDelete{func(app.DeleteSalesCommissionInput) error { return nil }},
		&mockCalc{},
	)
	err := adapter.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_Calculate_Success(t *testing.T) {
	calc := &app.CommissionCalculation{
		CommissionID: 1, InventoryID: 100, BasePriceUsed: 500, Quantity: 2,
		RatePercent: 10, CommissionAmt: 100, PriceSource: "base_price",
	}
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{},
		&mockCalc{func(app.CalculateCommissionInput) (*app.CommissionCalculation, error) { return calc, nil }},
	)
	resp, err := adapter.Calculate(salescommissioninterface.CalculateCommissionParams{InventoryID: 100, Quantity: 2})
	if err != nil {
		t.Fatal(err)
	}
	if resp.CommissionID != 1 || resp.CommissionAmt != 100 || resp.PriceSource != "base_price" {
		t.Error("unexpected calculation output")
	}
}

func TestAdapter_UpdateMaxPrice_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{func(int64, float64) error { return iface.ErrInternal }},
		&mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	err := adapter.UpdateMaxPrice(struct {
		CommissionID int64
		MaxPrice     float64
	}{CommissionID: 1, MaxPrice: 200})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
