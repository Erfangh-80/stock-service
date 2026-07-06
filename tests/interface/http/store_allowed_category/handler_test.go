package storeallowedcategoryhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	storeallowedcategory "stock-service/internal/domain/store_allowed_category"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	storeallowedcategoryinterface "stock-service/internal/interface/store_allowed_category"
)

type mockCreateCategory struct {
	fn func(storeID, categoryID int64) (*storeallowedcategory.StoreAllowedCategory, error)
}

func (m *mockCreateCategory) Execute(storeID, categoryID int64) (*storeallowedcategory.StoreAllowedCategory, error) {
	return m.fn(storeID, categoryID)
}

type mockApproveCategory struct {
	fn func(categoryID int64) error
}

func (m *mockApproveCategory) Execute(categoryID int64) error {
	return m.fn(categoryID)
}

type mockRejectCategory struct {
	fn func(categoryID int64) error
}

func (m *mockRejectCategory) Execute(categoryID int64) error {
	return m.fn(categoryID)
}

func TestStoreAllowedCategoryHandler_Create_Success(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{func(storeID, categoryID int64) (*storeallowedcategory.StoreAllowedCategory, error) {
			return &storeallowedcategory.StoreAllowedCategory{
				ID: 1, StoreID: storeID, CategoryID: categoryID,
				Status: storeallowedcategory.StatusPending,
			}, nil
		}},
		&mockApproveCategory{},
		&mockRejectCategory{},
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"store_id":1,"category_id":1}`
	req := httptest.NewRequest("POST", "/api/v1/store-categories", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var resp storeallowedcategoryinterface.StoreAllowedCategoryOutput
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.StoreID != 1 || resp.CategoryID != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestStoreAllowedCategoryHandler_Create_InvalidJSON(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{},
		&mockApproveCategory{},
		&mockRejectCategory{},
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/store-categories", strings.NewReader(`{invalid}`))
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

func TestStoreAllowedCategoryHandler_Approve_Success(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{},
		&mockApproveCategory{func(categoryID int64) error {
			return nil
		}},
		&mockRejectCategory{},
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/store-categories/1/approve", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestStoreAllowedCategoryHandler_Approve_InvalidID(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{},
		&mockApproveCategory{},
		&mockRejectCategory{},
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/store-categories/abc/approve", nil)
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

func TestStoreAllowedCategoryHandler_Reject_Success(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{},
		&mockApproveCategory{},
		&mockRejectCategory{func(categoryID int64) error {
			return nil
		}},
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/store-categories/1/reject", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestStoreAllowedCategoryHandler_Reject_InvalidID(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{},
		&mockApproveCategory{},
		&mockRejectCategory{},
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/api/v1/store-categories/abc/reject", nil)
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
