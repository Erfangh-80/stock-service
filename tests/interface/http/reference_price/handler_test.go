package referencepricehttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestReferencePriceHandler_Create_Success(t *testing.T) {
	adapter := referencepriceinterface.NewAdapter(
		&mockCreateReferencePrice{func(productID int32, price float64, source string) (*referenceprice.ReferencePrice, error) {
			return &referenceprice.ReferencePrice{ID: 1, ProductID: 42, Price: 99.99, Source: "manual"}, nil
		}},
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
	adapter := referencepriceinterface.NewAdapter(&mockCreateReferencePrice{})
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
	adapter := referencepriceinterface.NewAdapter(
		&mockCreateReferencePrice{func(productID int32, price float64, source string) (*referenceprice.ReferencePrice, error) {
			return nil, referenceprice.ErrInvalidReferencePrice
		}},
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
