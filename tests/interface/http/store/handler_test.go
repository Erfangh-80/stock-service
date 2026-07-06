package storehttp_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
	"stock-service/internal/interface/http/handler"
	storeinterface "stock-service/internal/interface/store"
)

type mockCreateStore struct{ fn func(appstore.CreateStoreInput) (*store.Store, error) }
type mockGetStore struct{ fn func(appstore.GetStoreInput) (*store.Store, error) }
type mockListStores struct{ fn func(appstore.ListStoresInput) (*appstore.ListStoresOutput, error) }
type mockToggleBulkSale struct{ fn func(appstore.ToggleBulkSaleInput) (*store.Store, error) }
type mockToggleCommission struct{ fn func(appstore.ToggleCommissionInput) (*store.Store, error) }
type mockUpdateContact struct{ fn func(appstore.UpdateContactInput) (*store.Store, error) }
type mockUpdateName struct{ fn func(appstore.UpdateStoreNameInput) (*store.Store, error) }
type mockUpdateProfile struct{ fn func(appstore.UpdateStoreProfileInput) (*store.Store, error) }
type mockDeleteStore struct{ fn func(appstore.DeleteStoreInput) error }

func (m *mockCreateStore) Execute(i appstore.CreateStoreInput) (*store.Store, error) { return m.fn(i) }
func (m *mockGetStore) Execute(i appstore.GetStoreInput) (*store.Store, error) { return m.fn(i) }
func (m *mockListStores) Execute(i appstore.ListStoresInput) (*appstore.ListStoresOutput, error) { return m.fn(i) }
func (m *mockToggleBulkSale) Execute(i appstore.ToggleBulkSaleInput) (*store.Store, error) { return m.fn(i) }
func (m *mockToggleCommission) Execute(i appstore.ToggleCommissionInput) (*store.Store, error) { return m.fn(i) }
func (m *mockUpdateContact) Execute(i appstore.UpdateContactInput) (*store.Store, error) { return m.fn(i) }
func (m *mockUpdateName) Execute(i appstore.UpdateStoreNameInput) (*store.Store, error) { return m.fn(i) }
func (m *mockUpdateProfile) Execute(i appstore.UpdateStoreProfileInput) (*store.Store, error) { return m.fn(i) }
func (m *mockDeleteStore) Execute(i appstore.DeleteStoreInput) error { return m.fn(i) }

func now() time.Time {
	t, _ := time.Parse(time.RFC3339, "2026-07-05T12:00:00Z")
	return t
}

func baseStore() *store.Store {
	return &store.Store{
		ID: 1, UserID: 10, StoreName: "Test Store",
		Status:                 store.StoreStatusActive,
		IsCommissionApplicable: true,
		IsBulkSaleEnabled:      false,
		CreatedAt:              now(),
	}
}

func TestHandler_Create_Success(t *testing.T) {
	s := baseStore()
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{func(i appstore.CreateStoreInput) (*store.Store, error) {
			if i.UserID != 10 || i.StoreName != "New Store" {
				t.Error("unexpected input")
			}
			return s, nil
		}},
		&mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/stores", strings.NewReader(`{"user_id":10,"store_name":"New Store"}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var resp storeinterface.StoreResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.StoreName != "Test Store" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestHandler_Create_InvalidInput(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{func(i appstore.CreateStoreInput) (*store.Store, error) {
			return nil, store.ErrStoreNameRequired
		}},
		&mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/stores", strings.NewReader(`{"user_id":10,"store_name":""}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var errResp struct{ Error string `json:"error"` }
	json.Unmarshal(b, &errResp)
	if errResp.Error != "invalid input" {
		t.Errorf("expected 'invalid input', got %q", errResp.Error)
	}
}

func TestHandler_List_Success(t *testing.T) {
	s := baseStore()
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil},
		&mockListStores{func(i appstore.ListStoresInput) (*appstore.ListStoresOutput, error) {
			return &appstore.ListStoresOutput{
				Stores: []*store.Store{s},
				Total:  1, Page: 1, Limit: 20,
			}, nil
		}},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/stores", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp storeinterface.ListStoresResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Stores) != 1 || resp.Total != 1 {
		t.Errorf("expected 1 store, got %d", len(resp.Stores))
	}
}

func TestHandler_List_WithFilter(t *testing.T) {
	s := baseStore()
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil},
		&mockListStores{func(i appstore.ListStoresInput) (*appstore.ListStoresOutput, error) {
			if i.UserID == nil || *i.UserID != 10 {
				t.Error("expected user_id filter")
			}
			return &appstore.ListStoresOutput{
				Stores: []*store.Store{s},
				Total:  1, Page: 1, Limit: 10,
			}, nil
		}},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/stores?user_id=10&limit=10", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	s := baseStore()
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil},
		&mockGetStore{func(i appstore.GetStoreInput) (*store.Store, error) {
			if i.ID != 1 {
				t.Error("unexpected id")
			}
			return s, nil
		}},
		&mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/stores/1", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp storeinterface.StoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestHandler_Get_NotFound(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil},
		&mockGetStore{func(i appstore.GetStoreInput) (*store.Store, error) {
			return nil, store.ErrStoreNotFound
		}},
		&mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/stores/999", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_Get_InvalidID(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/stores/abc", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var errResp struct{ Error string `json:"error"` }
	json.Unmarshal(b, &errResp)
	if errResp.Error != "invalid id" {
		t.Errorf("expected 'invalid id', got %q", errResp.Error)
	}
}

func TestHandler_UpdateName_Success(t *testing.T) {
	s := baseStore()
	s.StoreName = "Updated"
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{func(i appstore.UpdateStoreNameInput) (*store.Store, error) {
			if i.StoreID != 1 || i.Name != "Updated" {
				t.Error("unexpected input")
			}
			return s, nil
		}},
		&mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/stores/1", strings.NewReader(`{"name":"Updated"}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp storeinterface.StoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.StoreName != "Updated" {
		t.Errorf("expected 'Updated', got %q", resp.StoreName)
	}
}

func TestHandler_UpdateProfile_Success(t *testing.T) {
	addr := int64(99)
	s := baseStore()
	s.AddressID = &addr
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil},
		&mockUpdateProfile{func(i appstore.UpdateStoreProfileInput) (*store.Store, error) {
			if i.StoreID != 1 || *i.AddressID != 99 {
				t.Error("unexpected input")
			}
			return s, nil
		}},
		&mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/stores/1/profile", strings.NewReader(`{"address_id":99}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp storeinterface.StoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.AddressID == nil || *resp.AddressID != 99 {
		t.Errorf("expected address_id 99")
	}
}

func TestHandler_ToggleBulkSale_Success(t *testing.T) {
	s := baseStore()
	s.IsBulkSaleEnabled = true
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{func(i appstore.ToggleBulkSaleInput) (*store.Store, error) {
			if i.StoreID != 1 {
				t.Error("unexpected id")
			}
			return s, nil
		}},
		&mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/stores/1/bulk-sale", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp storeinterface.StoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.IsBulkSaleEnabled {
		t.Error("expected bulk sale enabled")
	}
}

func TestHandler_ToggleCommission_Success(t *testing.T) {
	s := baseStore()
	s.IsCommissionApplicable = false
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil},
		&mockToggleCommission{func(i appstore.ToggleCommissionInput) (*store.Store, error) {
			if i.StoreID != 1 {
				t.Error("unexpected id")
			}
			return s, nil
		}},
		&mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/stores/1/commission", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp storeinterface.StoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.IsCommissionApplicable {
		t.Error("expected commission not applicable")
	}
}

func TestHandler_UpdateContact_Success(t *testing.T) {
	phone := "+1234567890"
	s := baseStore()
	s.ContactPhone = &phone
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil},
		&mockUpdateContact{func(i appstore.UpdateContactInput) (*store.Store, error) {
			if i.StoreID != 1 || *i.ContactPhone != "+1234567890" {
				t.Error("unexpected input")
			}
			return s, nil
		}},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/stores/1/contact", strings.NewReader(`{"contact_phone":"+1234567890"}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp storeinterface.StoreResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ContactPhone == nil || *resp.ContactPhone != "+1234567890" {
		t.Error("expected contact phone to be set")
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil},
		&mockDeleteStore{func(i appstore.DeleteStoreInput) error { return nil }},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/api/v1/stores/1", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var resp map[string]string
	json.Unmarshal(b, &resp)
	if resp["status"] != "deleted" {
		t.Errorf("expected 'deleted', got %q", resp["status"])
	}
}

func TestHandler_InvalidJSON(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulkSale{nil}, &mockToggleCommission{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	mux := http.NewServeMux()
	handler.NewStoreHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/stores", strings.NewReader(`{invalid}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var errResp struct{ Error string `json:"error"` }
	json.Unmarshal(b, &errResp)
	if errResp.Error != "invalid JSON" {
		t.Errorf("expected 'invalid JSON', got %q", errResp.Error)
	}
}
