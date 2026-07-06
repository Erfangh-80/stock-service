package referencepricehttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	app "stock-service/internal/application/reference_price"
	referenceprice "stock-service/internal/domain/reference_price"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	referencepriceinterface "stock-service/internal/interface/reference_price"
)

type mockCreateReferencePrice struct {
	fn func(productID int32, price float64, source string) (*referenceprice.ReferencePrice, error)
}

func (m *mockCreateReferencePrice) Execute(productID int32, price float64, source string) (*referenceprice.ReferencePrice, error) {
	return m.fn(productID, price, source)
}

type mockGetReferencePrice struct {
	fn func(app.GetReferencePriceInput) (*referenceprice.ReferencePrice, error)
}

func (m *mockGetReferencePrice) Execute(input app.GetReferencePriceInput) (*referenceprice.ReferencePrice, error) {
	return m.fn(input)
}

type mockGetByProductReferencePrice struct {
	fn func(app.GetByProductReferencePriceInput) (*referenceprice.ReferencePrice, error)
}

func (m *mockGetByProductReferencePrice) Execute(input app.GetByProductReferencePriceInput) (*referenceprice.ReferencePrice, error) {
	return m.fn(input)
}

type mockListReferencePrices struct {
	fn func(app.ListReferencePricesInput) (*app.ListReferencePricesOutput, error)
}

func (m *mockListReferencePrices) Execute(input app.ListReferencePricesInput) (*app.ListReferencePricesOutput, error) {
	return m.fn(input)
}

type mockDeleteReferencePrice struct {
	fn func(app.DeleteReferencePriceInput) error
}

func (m *mockDeleteReferencePrice) Execute(input app.DeleteReferencePriceInput) error {
	return m.fn(input)
}

type mockValidateReferencePrice struct {
	fn func(app.ValidateReferencePriceInput) (*app.ReferencePriceValidation, error)
}

func (m *mockValidateReferencePrice) Execute(input app.ValidateReferencePriceInput) (*app.ReferencePriceValidation, error) {
	return m.fn(input)
}

func newTestAdapter(
	create *mockCreateReferencePrice,
	get *mockGetReferencePrice,
	getByProduct *mockGetByProductReferencePrice,
	list *mockListReferencePrices,
	del *mockDeleteReferencePrice,
	validate *mockValidateReferencePrice,
) *referencepriceinterface.Adapter {
	return referencepriceinterface.NewAdapter(create, get, getByProduct, list, del, validate)
}

func TestReferencePriceHandler_Create_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateReferencePrice{func(int32, float64, string) (*referenceprice.ReferencePrice, error) {
			return &referenceprice.ReferencePrice{ID: 1, ProductID: 42, Price: 99.99, Source: "manual"}, nil
		}},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"product_id":42,"price":99.99,"source":"manual"}`
	req := httptest.NewRequest("POST", "/api/v1/reference-prices", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var resp referencepriceinterface.ReferencePriceOutput
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.ProductID != 42 || resp.Price != 99.99 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestReferencePriceHandler_Create_InvalidJSON(t *testing.T) {
	adapter := newTestAdapter(&mockCreateReferencePrice{}, &mockGetReferencePrice{}, &mockGetByProductReferencePrice{}, &mockListReferencePrices{}, &mockDeleteReferencePrice{}, &mockValidateReferencePrice{})
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/reference-prices", strings.NewReader(`{invalid}`))
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

func TestReferencePriceHandler_Create_InvalidInput(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateReferencePrice{func(int32, float64, string) (*referenceprice.ReferencePrice, error) {
			return nil, referenceprice.ErrInvalidReferencePrice
		}},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"product_id":42,"price":0,"source":"test"}`
	req := httptest.NewRequest("POST", "/api/v1/reference-prices", strings.NewReader(body))
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

func TestReferencePriceHandler_Get_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{func(app.GetReferencePriceInput) (*referenceprice.ReferencePrice, error) {
			return &referenceprice.ReferencePrice{ID: 1, ProductID: 42, Price: 99.99, Source: "manual"}, nil
		}},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/reference-prices/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp referencepriceinterface.ReferencePriceOutput
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.ID != 1 || resp.ProductID != 42 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestReferencePriceHandler_Get_InvalidID(t *testing.T) {
	adapter := newTestAdapter(&mockCreateReferencePrice{}, &mockGetReferencePrice{}, &mockGetByProductReferencePrice{}, &mockListReferencePrices{}, &mockDeleteReferencePrice{}, &mockValidateReferencePrice{})
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/reference-prices/abc", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestReferencePriceHandler_Get_NotFound(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{func(app.GetReferencePriceInput) (*referenceprice.ReferencePrice, error) {
			return nil, referenceprice.ErrReferencePriceNotFound
		}},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/reference-prices/999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestReferencePriceHandler_GetByProduct_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{func(app.GetByProductReferencePriceInput) (*referenceprice.ReferencePrice, error) {
			return &referenceprice.ReferencePrice{ID: 1, ProductID: 42, Price: 99.99, Source: "manual"}, nil
		}},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/reference-prices/by-product/42", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp referencepriceinterface.ReferencePriceOutput
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.ProductID != 42 {
		t.Errorf("unexpected product_id: %d", resp.ProductID)
	}
}

func TestReferencePriceHandler_Delete_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{func(app.DeleteReferencePriceInput) error { return nil }},
		&mockValidateReferencePrice{},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("DELETE", "/api/v1/reference-prices/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestReferencePriceHandler_List_Success(t *testing.T) {
	output := &app.ListReferencePricesOutput{
		ReferencePrices: []*referenceprice.ReferencePrice{
			{ID: 1, ProductID: 42, Price: 99.99, Source: "manual"},
		},
		Total: 1, Page: 1, Limit: 20,
	}
	adapter := newTestAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{func(app.ListReferencePricesInput) (*app.ListReferencePricesOutput, error) { return output, nil }},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/reference-prices", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp referencepriceinterface.ReferencePriceListResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.ReferencePrices) != 1 || resp.Total != 1 {
		t.Errorf("unexpected list: %+v", resp)
	}
}

func TestReferencePriceHandler_Validate_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{func(app.ValidateReferencePriceInput) (*app.ReferencePriceValidation, error) {
			return &app.ReferencePriceValidation{
				ProductID: 42, ReferencePriceID: 1, ReferencePrice: 99.99,
				Source: "manual", InventoryCount: 1, BasePrices: []float64{100},
				Comparison: "within_range",
			}, nil
		}},
	)
	h := handler.NewReferencePriceHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/reference-prices/by-product/42/validate", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp referencepriceinterface.ValidationOutput
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.ProductID != 42 || resp.Comparison != "within_range" {
		t.Errorf("unexpected validation: %+v", resp)
	}
}
