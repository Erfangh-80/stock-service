package warehousehttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	warehousedomain "stock-service/internal/domain/warehouse"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	warehouseiface "stock-service/internal/interface/warehouse"
)

type mockCreateWarehouse struct {
	fn func(input warehouseiface.CreateWarehouseInput) (*warehousedomain.Warehouse, error)
}

func (m *mockCreateWarehouse) Execute(input warehouseiface.CreateWarehouseInput) (*warehousedomain.Warehouse, error) {
	return m.fn(input)
}

type mockUpdateVisibility struct {
	fn func(input warehouseiface.UpdateVisibilityInput) (*warehousedomain.Warehouse, error)
}

func (m *mockUpdateVisibility) Execute(input warehouseiface.UpdateVisibilityInput) (*warehousedomain.Warehouse, error) {
	return m.fn(input)
}

type mockUpdateContact struct {
	fn func(input warehouseiface.UpdateContactInput) (*warehousedomain.Warehouse, error)
}

func (m *mockUpdateContact) Execute(input warehouseiface.UpdateContactInput) (*warehousedomain.Warehouse, error) {
	return m.fn(input)
}

func TestWarehouseHandler_Create_Success(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{func(input warehouseiface.CreateWarehouseInput) (*warehousedomain.Warehouse, error) {
			return &warehousedomain.Warehouse{
				ID: 1, CreatedByUserID: input.CreatedByUserID, WarehouseName: input.WarehouseName,
				IsPublic: false, CollectionMethod: "pickup",
			}, nil
		}},
		&mockUpdateVisibility{},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"created_by_user_id":1,"warehouse_name":"main"}`
	req := httptest.NewRequest("POST", "/api/v1/warehouses", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var resp warehouseiface.WarehouseResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.CreatedByUserID != 1 || resp.WarehouseName != "main" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestWarehouseHandler_Create_InvalidJSON(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{},
		&mockUpdateVisibility{},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/warehouses", strings.NewReader(`{invalid}`))
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

func TestWarehouseHandler_Create_InvalidInput(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{func(input warehouseiface.CreateWarehouseInput) (*warehousedomain.Warehouse, error) {
			return nil, warehousedomain.ErrWarehouseNameRequired
		}},
		&mockUpdateVisibility{},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"created_by_user_id":1,"warehouse_name":""}`
	req := httptest.NewRequest("POST", "/api/v1/warehouses", strings.NewReader(body))
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

func TestWarehouseHandler_UpdateVisibility_Success(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{},
		&mockUpdateVisibility{func(input warehouseiface.UpdateVisibilityInput) (*warehousedomain.Warehouse, error) {
			return &warehousedomain.Warehouse{
				ID: input.WarehouseID, CreatedByUserID: 1, WarehouseName: "main",
				IsPublic: input.IsPublic, CollectionMethod: "pickup",
			}, nil
		}},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"is_public":true}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/visibility", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp warehouseiface.WarehouseResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || !resp.IsPublic {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestWarehouseHandler_UpdateVisibility_InvalidID(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{},
		&mockUpdateVisibility{},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"is_public":true}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/abc/visibility", strings.NewReader(body))
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

func TestWarehouseHandler_UpdateVisibility_InvalidJSON(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{},
		&mockUpdateVisibility{},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/visibility", strings.NewReader(`{invalid}`))
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

func TestWarehouseHandler_UpdateContact_Success(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{},
		&mockUpdateVisibility{},
		&mockUpdateContact{func(input warehouseiface.UpdateContactInput) (*warehousedomain.Warehouse, error) {
			return &warehousedomain.Warehouse{
				ID: input.WarehouseID, CreatedByUserID: 1, WarehouseName: "main",
				IsPublic: false, CollectionMethod: input.CollectionMethod,
				Phone: input.Phone, ContactPhone: input.ContactPhone,
			}, nil
		}},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"phone":"+123456789","collection_method":"pickup"}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/contact", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp warehouseiface.WarehouseResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.CollectionMethod != "pickup" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestWarehouseHandler_UpdateContact_InvalidID(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{},
		&mockUpdateVisibility{},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"phone":"+123456789","collection_method":"pickup"}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/abc/contact", strings.NewReader(body))
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

func TestWarehouseHandler_UpdateContact_InvalidJSON(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreateWarehouse{},
		&mockUpdateVisibility{},
		&mockUpdateContact{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/contact", strings.NewReader(`{invalid}`))
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
