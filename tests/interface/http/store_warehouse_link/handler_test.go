package storewarehouselinkhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	app "stock-service/internal/application/store_warehouse_link"
	domain "stock-service/internal/domain/store_warehouse_link"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	storewarehouselinkiface "stock-service/internal/interface/store_warehouse_link"
)

type mockCreateLink struct {
	fn func(int64, int64) (*domain.StoreWarehouseLink, error)
}

func (m *mockCreateLink) Execute(storeID, warehouseID int64) (*domain.StoreWarehouseLink, error) {
	return m.fn(storeID, warehouseID)
}

type mockGetLink struct {
	fn func(app.GetLinkInput) (*domain.StoreWarehouseLink, error)
}

func (m *mockGetLink) Execute(input app.GetLinkInput) (*domain.StoreWarehouseLink, error) {
	return m.fn(input)
}

type mockListLinks struct {
	fn func(app.ListLinksInput) (*app.ListLinksOutput, error)
}

func (m *mockListLinks) Execute(input app.ListLinksInput) (*app.ListLinksOutput, error) {
	return m.fn(input)
}

type mockChangeRelation struct {
	fn func(app.ChangeRelationInput) (*domain.StoreWarehouseLink, error)
}

func (m *mockChangeRelation) Execute(input app.ChangeRelationInput) (*domain.StoreWarehouseLink, error) {
	return m.fn(input)
}

type mockDeleteLink struct {
	fn func(app.DeleteLinkInput) error
}

func (m *mockDeleteLink) Execute(input app.DeleteLinkInput) error {
	return m.fn(input)
}

func newTestAdapter(
	create *mockCreateLink,
	get *mockGetLink,
	list *mockListLinks,
	change *mockChangeRelation,
	del *mockDeleteLink,
) *storewarehouselinkiface.Adapter {
	if create == nil {
		create = &mockCreateLink{}
	}
	if get == nil {
		get = &mockGetLink{}
	}
	if list == nil {
		list = &mockListLinks{}
	}
	if change == nil {
		change = &mockChangeRelation{}
	}
	if del == nil {
		del = &mockDeleteLink{}
	}
	return storewarehouselinkiface.NewAdapter(create, get, list, change, del)
}

func TestHandler_Create_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateLink{func(storeID, warehouseID int64) (*domain.StoreWarehouseLink, error) {
			return &domain.StoreWarehouseLink{
				ID: 1, StoreID: storeID, WarehouseID: warehouseID,
				RelationType: domain.RelationTypePrimary,
			}, nil
		}},
		nil, nil, nil, nil,
	)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"store_id":1,"warehouse_id":1}`
	req := httptest.NewRequest("POST", "/api/v1/warehouse-links", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var resp storewarehouselinkiface.LinkResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.StoreID != 1 || resp.WarehouseID != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/warehouse-links", strings.NewReader(`{invalid}`))
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

func TestHandler_Get_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil,
		&mockGetLink{func(input app.GetLinkInput) (*domain.StoreWarehouseLink, error) {
			return &domain.StoreWarehouseLink{
				ID: 1, StoreID: 10, WarehouseID: 20, RelationType: domain.RelationTypePrimary,
			}, nil
		}},
		nil, nil, nil,
	)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/warehouse-links/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp storewarehouselinkiface.LinkResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestHandler_Get_InvalidID(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/warehouse-links/abc", nil)
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

func TestHandler_List_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil,
		&mockListLinks{func(input app.ListLinksInput) (*app.ListLinksOutput, error) {
			return &app.ListLinksOutput{
				Links: []*domain.StoreWarehouseLink{
					{ID: 1, StoreID: 10, WarehouseID: 20, RelationType: domain.RelationTypePrimary},
				},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil,
	)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/warehouse-links", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp storewarehouselinkiface.LinkListResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.Links) != 1 || resp.Total != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestHandler_List_FilterByStore(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil,
		&mockListLinks{func(input app.ListLinksInput) (*app.ListLinksOutput, error) {
			if input.StoreID == nil || *input.StoreID != 1 {
				t.Error("expected store_id=1 filter")
			}
			return &app.ListLinksOutput{
				Links: []*domain.StoreWarehouseLink{
					{ID: 1, StoreID: 1, WarehouseID: 10, RelationType: domain.RelationTypePrimary},
				},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil,
	)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/warehouse-links?store_id=1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil, nil,
		&mockDeleteLink{func(input app.DeleteLinkInput) error {
			return nil
		}},
	)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("DELETE", "/api/v1/warehouse-links/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("DELETE", "/api/v1/warehouse-links/abc", nil)
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

func TestHandler_ChangeRelation_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil,
		&mockChangeRelation{func(input app.ChangeRelationInput) (*domain.StoreWarehouseLink, error) {
			return &domain.StoreWarehouseLink{
				ID: input.LinkID, StoreID: 1, WarehouseID: 1,
				RelationType: domain.RelationType(input.RelationType),
			}, nil
		}},
		nil,
	)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"relation_type":"primary"}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouse-links/1/relation", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp storewarehouselinkiface.LinkResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestHandler_ChangeRelation_InvalidID(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"relation_type":"primary"}`
	req := httptest.NewRequest("PUT", "/api/v1/warehouse-links/abc/relation", strings.NewReader(body))
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

func TestHandler_ChangeRelation_InvalidJSON(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil)
	h := handler.NewStoreWarehouseLinkHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("PUT", "/api/v1/warehouse-links/1/relation", strings.NewReader(`{invalid}`))
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
