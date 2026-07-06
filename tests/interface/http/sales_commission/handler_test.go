package salescommissionhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	app "stock-service/internal/application/sales_commission"
	salescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/interface/http/handler"
	salescommissioninterface "stock-service/internal/interface/sales_commission"
)

type mockCreate struct {
	fn func(int64, int64, salescommission.SaleModel, float64, float64) (*salescommission.SalesCommission, error)
}
func (m *mockCreate) Execute(a, b int64, c salescommission.SaleModel, d, e float64) (*salescommission.SalesCommission, error) { return m.fn(a, b, c, d, e) }

type mockUpdateMaxPrice struct {
	fn func(int64, float64) error
}
func (m *mockUpdateMaxPrice) Execute(a int64, b float64) error { return m.fn(a, b) }

type mockUpdateMinQty struct {
	fn func(int64, int) error
}
func (m *mockUpdateMinQty) Execute(a int64, b int) error { return m.fn(a, b) }

type mockGet struct {
	fn func(app.GetSalesCommissionInput) (*salescommission.SalesCommission, error)
}
func (m *mockGet) Execute(a app.GetSalesCommissionInput) (*salescommission.SalesCommission, error) { return m.fn(a) }

type mockGetByInv struct {
	fn func(app.GetByInventorySalesCommissionInput) (*salescommission.SalesCommission, error)
}
func (m *mockGetByInv) Execute(a app.GetByInventorySalesCommissionInput) (*salescommission.SalesCommission, error) { return m.fn(a) }

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

func TestSalesCommissionHandler_Create_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{func(int64, int64, salescommission.SaleModel, float64, float64) (*salescommission.SalesCommission, error) {
			return &salescommission.SalesCommission{ID: 1, InventoryID: 1, RatePercent: 10, MinPrice: 100}, nil
		}},
		&mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"inventory_id":1,"category_commission_rule_id":1,"sale_model":"retail","rate_percent":10,"min_price":100}`
	req := httptest.NewRequest("POST", "/api/v1/sales-commissions", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var resp salescommissioninterface.SalesCommissionOutput
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.RatePercent != 10 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestSalesCommissionHandler_Create_InvalidJSON(t *testing.T) {
	adapter := newTestAdapter(&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{})
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/sales-commissions", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_Get_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{},
		&mockGet{func(app.GetSalesCommissionInput) (*salescommission.SalesCommission, error) {
			return &salescommission.SalesCommission{ID: 1, InventoryID: 100, RatePercent: 10, MinPrice: 50}, nil
		}},
		&mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/sales-commissions/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_Get_NotFound(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{},
		&mockGet{func(app.GetSalesCommissionInput) (*salescommission.SalesCommission, error) { return nil, salescommission.ErrCommissionNotFound }},
		&mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/sales-commissions/999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_GetByInventory_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{},
		&mockGetByInv{func(app.GetByInventorySalesCommissionInput) (*salescommission.SalesCommission, error) {
			return &salescommission.SalesCommission{ID: 1, InventoryID: 100, RatePercent: 10, MinPrice: 50}, nil
		}},
		&mockList{}, &mockDelete{}, &mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/sales-commissions/by-inventory/100", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_Delete_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{},
		&mockDelete{func(app.DeleteSalesCommissionInput) error { return nil }},
		&mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("DELETE", "/api/v1/sales-commissions/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_List_Success(t *testing.T) {
	output := &app.ListSalesCommissionsOutput{
		Commissions: []*salescommission.SalesCommission{
			{ID: 1, InventoryID: 100, SaleModel: "retail", RatePercent: 10, MinPrice: 50},
		},
		Total: 1, Page: 1, Limit: 20,
	}
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{},
		&mockList{func(app.ListSalesCommissionsInput) (*app.ListSalesCommissionsOutput, error) { return output, nil }},
		&mockDelete{}, &mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/sales-commissions", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp salescommissioninterface.SalesCommissionListResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.Commissions) != 1 {
		t.Errorf("expected 1 commission, got %d", len(resp.Commissions))
	}
}

func TestSalesCommissionHandler_Calculate_Success(t *testing.T) {
	calc := &app.CommissionCalculation{
		CommissionID: 1, InventoryID: 100, BasePriceUsed: 500, Quantity: 2,
		RatePercent: 10, CommissionAmt: 100, PriceSource: "base_price",
	}
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{},
		&mockCalc{func(app.CalculateCommissionInput) (*app.CommissionCalculation, error) { return calc, nil }},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"inventory_id":100,"quantity":2}`
	req := httptest.NewRequest("POST", "/api/v1/sales-commissions/calculate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp salescommissioninterface.CommissionCalculationOutput
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.CommissionAmt != 100 || resp.PriceSource != "base_price" {
		t.Errorf("unexpected calculation: %+v", resp)
	}
}

func TestSalesCommissionHandler_UpdateMaxPrice_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{func(int64, float64) error { return nil }},
		&mockUpdateMinQty{}, &mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/1/max-price", strings.NewReader(`{"max_price":200}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_UpdateMinQty_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreate{}, &mockUpdateMaxPrice{}, &mockUpdateMinQty{func(int64, int) error { return nil }},
		&mockGet{}, &mockGetByInv{}, &mockList{}, &mockDelete{}, &mockCalc{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/1/min-qty", strings.NewReader(`{"min_qty":5}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
