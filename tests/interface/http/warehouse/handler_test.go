package warehousehttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	app "stock-service/internal/application/warehouse"
	warehousedomain "stock-service/internal/domain/warehouse"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	warehouseiface "stock-service/internal/interface/warehouse"
)

type mockCreate struct {
	fn func(int64, string) (*warehousedomain.Warehouse, error)
}

func (m *mockCreate) Execute(uid int64, name string) (*warehousedomain.Warehouse, error) {
	return m.fn(uid, name)
}

type mockGet struct {
	fn func(int64) (*warehousedomain.Warehouse, error)
}

func (m *mockGet) Execute(id int64) (*warehousedomain.Warehouse, error) {
	if m.fn != nil {
		return m.fn(id)
	}
	return &warehousedomain.Warehouse{}, nil
}

type mockList struct{}

func (m *mockList) Execute(input app.ListWarehousesInput) (*app.ListWarehousesOutput, error) {
	return &app.ListWarehousesOutput{}, nil
}

type mockDel struct{}

func (m *mockDel) Execute(id int64) error { return nil }

type mockVis struct{}

func (m *mockVis) Execute(id int64, isPublic bool) error { return nil }

type mockCont struct{}

func (m *mockCont) Execute(id int64, phone, contactPhone *string, collectionMethod string) error {
	return nil
}

type mockUpd struct{}

func (m *mockUpd) Execute(input app.UpdateWarehouseInput) (*warehousedomain.Warehouse, error) {
	return &warehousedomain.Warehouse{
		ID:               input.WarehouseID,
		CreatedByUserID:  1,
		WarehouseName:    "main",
		IsPublic:         false,
		CollectionMethod: "pickup",
		CreatedAt:        time.Now(),
	}, nil
}

func baseHandler() *handler.WarehouseHandler {
	return handler.NewWarehouseHandler(
		warehouseiface.NewAdapter(
			&mockCreate{},
			&mockGet{},
			&mockList{},
			&mockDel{},
			&mockVis{},
			&mockCont{},
			&mockUpd{},
		),
	)
}

func register(h *handler.WarehouseHandler) *http.ServeMux {
	mux := http.NewServeMux()
	h.Register(mux)
	return mux
}

func TestWarehouseHandler_Create_Success(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreate{fn: func(uid int64, name string) (*warehousedomain.Warehouse, error) {
			return &warehousedomain.Warehouse{
				ID: 1, CreatedByUserID: uid, WarehouseName: name,
				IsPublic: false, CollectionMethod: "pickup",
				CreatedAt: time.Now(),
			}, nil
		}},
		&mockGet{}, &mockList{}, &mockDel{}, &mockVis{}, &mockCont{}, &mockUpd{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := register(h)

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
	h := baseHandler()
	mux := register(h)

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
		&mockCreate{fn: func(uid int64, name string) (*warehousedomain.Warehouse, error) {
			return nil, warehousedomain.ErrWarehouseNameRequired
		}},
		&mockGet{}, &mockList{}, &mockDel{}, &mockVis{}, &mockCont{}, &mockUpd{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := register(h)

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

func TestWarehouseHandler_Get_Success(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreate{},
		&mockGet{fn: func(id int64) (*warehousedomain.Warehouse, error) {
			return &warehousedomain.Warehouse{
				ID: id, CreatedByUserID: 1, WarehouseName: "main",
				IsPublic: false, CollectionMethod: "pickup", CreatedAt: time.Now(),
			}, nil
		}},
		&mockList{}, &mockDel{}, &mockVis{}, &mockCont{}, &mockUpd{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := register(h)

	req := httptest.NewRequest("GET", "/api/v1/warehouses/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp warehouseiface.WarehouseResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.WarehouseName != "main" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestWarehouseHandler_Get_InvalidID(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	req := httptest.NewRequest("GET", "/api/v1/warehouses/abc", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestWarehouseHandler_Get_NotFound(t *testing.T) {
	adapter := warehouseiface.NewAdapter(
		&mockCreate{},
		&mockGet{fn: func(id int64) (*warehousedomain.Warehouse, error) {
			return nil, warehousedomain.ErrWarehouseNotFound
		}},
		&mockList{}, &mockDel{}, &mockVis{}, &mockCont{}, &mockUpd{},
	)
	h := handler.NewWarehouseHandler(adapter)
	mux := register(h)

	req := httptest.NewRequest("GET", "/api/v1/warehouses/999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestWarehouseHandler_List_Success(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	req := httptest.NewRequest("GET", "/api/v1/warehouses", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestWarehouseHandler_Update_Success(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	body := `{"warehouse_name":"Updated Name"}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestWarehouseHandler_Delete_Success(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	req := httptest.NewRequest("DELETE", "/api/v1/warehouses/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}

func TestWarehouseHandler_UpdateVisibility_Success(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	body := `{"is_public":true}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/visibility", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestWarehouseHandler_UpdateVisibility_InvalidID(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	body := `{"is_public":true}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/abc/visibility", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestWarehouseHandler_UpdateVisibility_InvalidJSON(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/visibility", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestWarehouseHandler_UpdateContact_Success(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	body := `{"phone":"+123456789","collection_method":"pickup"}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/contact", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestWarehouseHandler_UpdateContact_InvalidID(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	body := `{"phone":"+123456789","collection_method":"pickup"}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouses/abc/contact", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestWarehouseHandler_UpdateContact_InvalidJSON(t *testing.T) {
	h := baseHandler()
	mux := register(h)

	req := httptest.NewRequest("PUT", "/api/v1/warehouses/1/contact", strings.NewReader(`{invalid}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
