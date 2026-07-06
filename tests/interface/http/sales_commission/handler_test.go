package salescommissionhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	salescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	salescommissioninterface "stock-service/internal/interface/sales_commission"
)

type mockCreateSalesCommission struct {
	fn func(inventoryID, categoryCommissionRuleID int64, saleModel salescommission.SaleModel, ratePercent, minPrice float64) (*salescommission.SalesCommission, error)
}

func (m *mockCreateSalesCommission) Execute(inventoryID, categoryCommissionRuleID int64, saleModel salescommission.SaleModel, ratePercent, minPrice float64) (*salescommission.SalesCommission, error) {
	return m.fn(inventoryID, categoryCommissionRuleID, saleModel, ratePercent, minPrice)
}

type mockUpdateMaxPrice struct {
	fn func(commissionID int64, maxPrice float64) error
}

func (m *mockUpdateMaxPrice) Execute(commissionID int64, maxPrice float64) error {
	return m.fn(commissionID, maxPrice)
}

type mockUpdateMinQty struct {
	fn func(commissionID int64, minQty int) error
}

func (m *mockUpdateMinQty) Execute(commissionID int64, minQty int) error {
	return m.fn(commissionID, minQty)
}

func TestSalesCommissionHandler_Create_Success(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{func(inventoryID, categoryCommissionRuleID int64, saleModel salescommission.SaleModel, ratePercent, minPrice float64) (*salescommission.SalesCommission, error) {
			return &salescommission.SalesCommission{
				ID: 1, InventoryID: inventoryID, CategoryCommissionRuleID: categoryCommissionRuleID,
				SaleModel: saleModel, RatePercent: ratePercent, MinPrice: minPrice,
			}, nil
		}},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{},
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
	if resp.ID != 1 || resp.InventoryID != 1 || resp.RatePercent != 10 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestSalesCommissionHandler_Create_InvalidJSON(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/sales-commissions", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp dto.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "invalid JSON" {
		t.Errorf("expected 'invalid JSON', got %q", errResp.Error)
	}
}

func TestSalesCommissionHandler_Create_InvalidInput(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{func(inventoryID, categoryCommissionRuleID int64, saleModel salescommission.SaleModel, ratePercent, minPrice float64) (*salescommission.SalesCommission, error) {
			return nil, salescommission.ErrInvalidRatePercent
		}},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"inventory_id":1,"category_commission_rule_id":1,"sale_model":"retail","rate_percent":150,"min_price":100}`
	req := httptest.NewRequest("POST", "/api/v1/sales-commissions", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp dto.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "invalid input" {
		t.Errorf("expected 'invalid input', got %q", errResp.Error)
	}
}

func TestSalesCommissionHandler_UpdateMaxPrice_Success(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{},
		&mockUpdateMaxPrice{func(commissionID int64, maxPrice float64) error {
			return nil
		}},
		&mockUpdateMinQty{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"max_price":200}`
	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/1/max-price", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_UpdateMaxPrice_InvalidID(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"max_price":200}`
	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/abc/max-price", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp dto.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "invalid id" {
		t.Errorf("expected 'invalid id', got %q", errResp.Error)
	}
}

func TestSalesCommissionHandler_UpdateMaxPrice_InvalidJSON(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/1/max-price", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp dto.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "invalid JSON" {
		t.Errorf("expected 'invalid JSON', got %q", errResp.Error)
	}
}

func TestSalesCommissionHandler_UpdateMinQty_Success(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{func(commissionID int64, minQty int) error {
			return nil
		}},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"min_qty":5}`
	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/1/min-qty", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSalesCommissionHandler_UpdateMinQty_InvalidID(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"min_qty":5}`
	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/abc/min-qty", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp dto.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "invalid id" {
		t.Errorf("expected 'invalid id', got %q", errResp.Error)
	}
}

func TestSalesCommissionHandler_UpdateMinQty_InvalidJSON(t *testing.T) {
	adapter := salescommissioninterface.NewAdapter(
		&mockCreateSalesCommission{},
		&mockUpdateMaxPrice{},
		&mockUpdateMinQty{},
	)
	h := handler.NewSalesCommissionHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/api/v1/sales-commissions/1/min-qty", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	var errResp dto.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "invalid JSON" {
		t.Errorf("expected 'invalid JSON', got %q", errResp.Error)
	}
}
