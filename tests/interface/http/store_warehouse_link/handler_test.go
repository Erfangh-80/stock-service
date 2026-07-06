package storewarehouselinkhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	storewarehouselinkdomain "stock-service/internal/domain/store_warehouse_link"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	storewarehouselinkiface "stock-service/internal/interface/store_warehouse_link"
)

type mockCreateLink struct {
	fn func(input storewarehouselinkiface.CreateLinkInput) (*storewarehouselinkdomain.StoreWarehouseLink, error)
}

func (m *mockCreateLink) Execute(input storewarehouselinkiface.CreateLinkInput) (*storewarehouselinkdomain.StoreWarehouseLink, error) {
	return m.fn(input)
}

type mockChangeRelation struct {
	fn func(input storewarehouselinkiface.ChangeRelationInput) (*storewarehouselinkdomain.StoreWarehouseLink, error)
}

func (m *mockChangeRelation) Execute(input storewarehouselinkiface.ChangeRelationInput) (*storewarehouselinkdomain.StoreWarehouseLink, error) {
	return m.fn(input)
}

func TestStoreWarehouseLinkHandler_Create_Success(t *testing.T) {
	adapter := storewarehouselinkiface.NewAdapter(
		&mockCreateLink{func(input storewarehouselinkiface.CreateLinkInput) (*storewarehouselinkdomain.StoreWarehouseLink, error) {
			return &storewarehouselinkdomain.StoreWarehouseLink{
				ID: 1, StoreID: input.StoreID, WarehouseID: input.WarehouseID,
				RelationType: storewarehouselinkdomain.RelationTypePrimary,
			}, nil
		}},
		&mockChangeRelation{},
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

func TestStoreWarehouseLinkHandler_Create_InvalidJSON(t *testing.T) {
	adapter := storewarehouselinkiface.NewAdapter(
		&mockCreateLink{},
		&mockChangeRelation{},
	)
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

func TestStoreWarehouseLinkHandler_ChangeRelation_Success(t *testing.T) {
	adapter := storewarehouselinkiface.NewAdapter(
		&mockCreateLink{},
		&mockChangeRelation{func(input storewarehouselinkiface.ChangeRelationInput) (*storewarehouselinkdomain.StoreWarehouseLink, error) {
			return &storewarehouselinkdomain.StoreWarehouseLink{
				ID: input.LinkID, StoreID: 1, WarehouseID: 1,
				RelationType: storewarehouselinkdomain.RelationType(input.RelationType),
			}, nil
		}},
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

func TestStoreWarehouseLinkHandler_ChangeRelation_InvalidID(t *testing.T) {
	adapter := storewarehouselinkiface.NewAdapter(
		&mockCreateLink{},
		&mockChangeRelation{},
	)
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

func TestStoreWarehouseLinkHandler_ChangeRelation_InvalidJSON(t *testing.T) {
	adapter := storewarehouselinkiface.NewAdapter(
		&mockCreateLink{},
		&mockChangeRelation{},
	)
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
