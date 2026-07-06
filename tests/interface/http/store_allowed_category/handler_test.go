package storeallowedcategoryhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	app "stock-service/internal/application/store_allowed_category"
	storeallowedcategory "stock-service/internal/domain/store_allowed_category"
	"stock-service/internal/interface/http/handler"
	"stock-service/internal/interface/http/dto"
	storeallowedcategoryinterface "stock-service/internal/interface/store_allowed_category"
)

type mockCreateCategory struct {
	fn func(int64, int64) (*storeallowedcategory.StoreAllowedCategory, error)
}

func (m *mockCreateCategory) Execute(storeID, categoryID int64) (*storeallowedcategory.StoreAllowedCategory, error) {
	return m.fn(storeID, categoryID)
}

type mockGetCategory struct {
	fn func(app.GetStoreCategoryInput) (*storeallowedcategory.StoreAllowedCategory, error)
}

func (m *mockGetCategory) Execute(input app.GetStoreCategoryInput) (*storeallowedcategory.StoreAllowedCategory, error) {
	return m.fn(input)
}

type mockListCategories struct {
	fn func(app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error)
}

func (m *mockListCategories) Execute(input app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error) {
	return m.fn(input)
}

type mockApproveCategory struct {
	fn func(app.ApproveCategoryInput) error
}

func (m *mockApproveCategory) Execute(input app.ApproveCategoryInput) error {
	return m.fn(input)
}

type mockRejectCategory struct {
	fn func(app.RejectCategoryInput) error
}

func (m *mockRejectCategory) Execute(input app.RejectCategoryInput) error {
	return m.fn(input)
}

type mockDeleteCategory struct {
	fn func(app.DeleteStoreCategoryInput) error
}

func (m *mockDeleteCategory) Execute(input app.DeleteStoreCategoryInput) error {
	return m.fn(input)
}

type mockValidateCategory struct {
	fn func(app.ValidateCategoryExistsInput) error
}

func (m *mockValidateCategory) Execute(input app.ValidateCategoryExistsInput) error {
	return m.fn(input)
}

func newTestAdapter(
	create *mockCreateCategory,
	get *mockGetCategory,
	list *mockListCategories,
	approve *mockApproveCategory,
	reject *mockRejectCategory,
	del *mockDeleteCategory,
	validate *mockValidateCategory,
) *storeallowedcategoryinterface.Adapter {
	if create == nil {
		create = &mockCreateCategory{}
	}
	if get == nil {
		get = &mockGetCategory{}
	}
	if list == nil {
		list = &mockListCategories{}
	}
	if approve == nil {
		approve = &mockApproveCategory{}
	}
	if reject == nil {
		reject = &mockRejectCategory{}
	}
	if del == nil {
		del = &mockDeleteCategory{}
	}
	if validate == nil {
		validate = &mockValidateCategory{}
	}
	return storeallowedcategoryinterface.NewAdapter(create, get, list, approve, reject, del, validate)
}

func TestStoreAllowedCategoryHandler_Create_Success(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateCategory{func(storeID, categoryID int64) (*storeallowedcategory.StoreAllowedCategory, error) {
			return &storeallowedcategory.StoreAllowedCategory{
				ID: 1, StoreID: storeID, CategoryID: categoryID,
				Status: storeallowedcategory.StatusPending,
			}, nil
		}},
		nil, nil, nil, nil, nil,
		&mockValidateCategory{func(app.ValidateCategoryExistsInput) error { return nil }},
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
	adapter := newTestAdapter(nil, nil, nil, nil, nil, nil, nil)
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

func TestStoreAllowedCategoryHandler_Get_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil,
		&mockGetCategory{func(input app.GetStoreCategoryInput) (*storeallowedcategory.StoreAllowedCategory, error) {
			return &storeallowedcategory.StoreAllowedCategory{
				ID: 1, StoreID: 100, CategoryID: 200, Status: storeallowedcategory.StatusPending,
			}, nil
		}},
		nil, nil, nil, nil, nil,
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/store-categories/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp storeallowedcategoryinterface.StoreAllowedCategoryOutput
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestStoreAllowedCategoryHandler_Get_InvalidID(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil, nil, nil)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/store-categories/abc", nil)
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

func TestStoreAllowedCategoryHandler_List_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil,
		&mockListCategories{func(input app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error) {
			return &app.ListStoreCategoriesOutput{
				Categories: []*storeallowedcategory.StoreAllowedCategory{
					{ID: 1, StoreID: 100, CategoryID: 200, Status: storeallowedcategory.StatusPending},
				},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil, nil, nil,
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/store-categories", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp storeallowedcategoryinterface.StoreCategoryListResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.Categories) != 1 || resp.Total != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestStoreAllowedCategoryHandler_List_FilterByStore(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil,
		&mockListCategories{func(input app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error) {
			if input.StoreID == nil || *input.StoreID != 1 {
				t.Error("expected store_id=1 filter")
			}
			return &app.ListStoreCategoriesOutput{
				Categories: []*storeallowedcategory.StoreAllowedCategory{
					{ID: 1, StoreID: 1, CategoryID: 200, Status: storeallowedcategory.StatusPending},
				},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil, nil, nil,
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/api/v1/store-categories?store_id=1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestStoreAllowedCategoryHandler_Delete_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil, nil, nil,
		&mockDeleteCategory{func(input app.DeleteStoreCategoryInput) error {
			return nil
		}},
		nil,
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("DELETE", "/api/v1/store-categories/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestStoreAllowedCategoryHandler_Delete_InvalidID(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil, nil, nil)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("DELETE", "/api/v1/store-categories/abc", nil)
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

func TestStoreAllowedCategoryHandler_Approve_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil,
		&mockApproveCategory{func(input app.ApproveCategoryInput) error {
			return nil
		}},
		nil, nil, nil,
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
	adapter := newTestAdapter(nil, nil, nil, nil, nil, nil, nil)
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
	adapter := newTestAdapter(
		nil, nil, nil, nil,
		&mockRejectCategory{func(input app.RejectCategoryInput) error {
			if input.CategoryID != 1 || input.SupportNote != "bad category" {
				t.Error("unexpected input")
			}
			return nil
		}},
		nil, nil,
	)
	h := handler.NewStoreAllowedCategoryHandler(adapter)
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"support_note":"bad category"}`
	req := httptest.NewRequest("POST", "/api/v1/store-categories/1/reject", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestStoreAllowedCategoryHandler_Reject_InvalidID(t *testing.T) {
	adapter := newTestAdapter(nil, nil, nil, nil, nil, nil, nil)
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
