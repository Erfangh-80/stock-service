package storeinterface_test

import (
	"testing"
	"time"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
	iface "stock-service/internal/interface"
	storeinterface "stock-service/internal/interface/store"
)

type mockCreateStore struct{ fn func(appstore.CreateStoreInput) (*store.Store, error) }
type mockGetStore struct{ fn func(appstore.GetStoreInput) (*store.Store, error) }
type mockListStores struct{ fn func(appstore.ListStoresInput) (*appstore.ListStoresOutput, error) }
type mockToggleBulk struct{ fn func(appstore.ToggleBulkSaleInput) (*store.Store, error) }
type mockToggleComm struct{ fn func(appstore.ToggleCommissionInput) (*store.Store, error) }
type mockUpdateContact struct{ fn func(appstore.UpdateContactInput) (*store.Store, error) }
type mockUpdateName struct{ fn func(appstore.UpdateStoreNameInput) (*store.Store, error) }
type mockUpdateProfile struct{ fn func(appstore.UpdateStoreProfileInput) (*store.Store, error) }
type mockDeleteStore struct{ fn func(appstore.DeleteStoreInput) error }

func (m *mockCreateStore) Execute(i appstore.CreateStoreInput) (*store.Store, error) { return m.fn(i) }
func (m *mockGetStore) Execute(i appstore.GetStoreInput) (*store.Store, error) { return m.fn(i) }
func (m *mockListStores) Execute(i appstore.ListStoresInput) (*appstore.ListStoresOutput, error) { return m.fn(i) }
func (m *mockToggleBulk) Execute(i appstore.ToggleBulkSaleInput) (*store.Store, error) { return m.fn(i) }
func (m *mockToggleComm) Execute(i appstore.ToggleCommissionInput) (*store.Store, error) { return m.fn(i) }
func (m *mockUpdateContact) Execute(i appstore.UpdateContactInput) (*store.Store, error) { return m.fn(i) }
func (m *mockUpdateName) Execute(i appstore.UpdateStoreNameInput) (*store.Store, error) { return m.fn(i) }
func (m *mockUpdateProfile) Execute(i appstore.UpdateStoreProfileInput) (*store.Store, error) { return m.fn(i) }
func (m *mockDeleteStore) Execute(i appstore.DeleteStoreInput) error { return m.fn(i) }

func baseStore() *store.Store {
	return &store.Store{
		ID: 1, UserID: 1, StoreName: "My Store",
		Status: store.StoreStatusActive,
		IsCommissionApplicable: true, IsBulkSaleEnabled: false,
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestStoreAdapter_Create_Success(t *testing.T) {
	s := baseStore()
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{func(i appstore.CreateStoreInput) (*store.Store, error) { return s, nil }},
		&mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	resp, err := adapter.Create(storeinterface.CreateStoreParams{UserID: 1, StoreName: "My Store"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.StoreName != "My Store" {
		t.Error("unexpected response")
	}
}

func TestStoreAdapter_Create_InvalidInput(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{func(i appstore.CreateStoreInput) (*store.Store, error) {
			return nil, store.ErrStoreNameRequired
		}},
		&mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	_, err := adapter.Create(storeinterface.CreateStoreParams{UserID: 1, StoreName: ""})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestStoreAdapter_Get_Success(t *testing.T) {
	s := baseStore()
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil},
		&mockGetStore{func(i appstore.GetStoreInput) (*store.Store, error) { return s, nil }},
		&mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	resp, err := adapter.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("expected ID 1")
	}
}

func TestStoreAdapter_Get_NotFound(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil},
		&mockGetStore{func(i appstore.GetStoreInput) (*store.Store, error) {
			return nil, store.ErrStoreNotFound
		}},
		&mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	_, err := adapter.Get(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStoreAdapter_List_Success(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil},
		&mockListStores{func(i appstore.ListStoresInput) (*appstore.ListStoresOutput, error) {
			return &appstore.ListStoresOutput{Stores: []*store.Store{baseStore()}, Total: 1, Page: 1, Limit: 20}, nil
		}},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	resp, err := adapter.List(storeinterface.ListStoresFilter{Page: 1, Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Stores) != 1 || resp.Total != 1 {
		t.Error("unexpected list response")
	}
}

func TestStoreAdapter_ToggleBulkSale_Success(t *testing.T) {
	s := baseStore()
	s.IsBulkSaleEnabled = true
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{func(i appstore.ToggleBulkSaleInput) (*store.Store, error) { return s, nil }},
		&mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	resp, err := adapter.ToggleBulkSale(1)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsBulkSaleEnabled {
		t.Error("expected bulk sale enabled")
	}
}

func TestStoreAdapter_ToggleCommission_Success(t *testing.T) {
	s := baseStore()
	s.IsCommissionApplicable = false
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil},
		&mockToggleComm{func(i appstore.ToggleCommissionInput) (*store.Store, error) { return s, nil }},
		&mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	resp, err := adapter.ToggleCommission(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsCommissionApplicable {
		t.Error("expected commission disabled")
	}
}

func TestStoreAdapter_UpdateContact_Success(t *testing.T) {
	phone := "+123"
	s := baseStore()
	s.ContactPhone = &phone
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil},
		&mockUpdateContact{func(i appstore.UpdateContactInput) (*store.Store, error) { return s, nil }},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	resp, err := adapter.UpdateContact(storeinterface.UpdateContactParams{StoreID: 1, ContactPhone: &phone})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ContactPhone == nil || *resp.ContactPhone != "+123" {
		t.Error("unexpected phone")
	}
}

func TestStoreAdapter_UpdateName_Success(t *testing.T) {
	s := baseStore()
	s.StoreName = "New Name"
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{func(i appstore.UpdateStoreNameInput) (*store.Store, error) { return s, nil }},
		&mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	resp, err := adapter.UpdateName(storeinterface.UpdateNameParams{StoreID: 1, Name: "New Name"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StoreName != "New Name" {
		t.Errorf("expected 'New Name', got %q", resp.StoreName)
	}
}

func TestStoreAdapter_UpdateProfile_Success(t *testing.T) {
	addr := int64(99)
	s := baseStore()
	s.AddressID = &addr
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil},
		&mockUpdateProfile{func(i appstore.UpdateStoreProfileInput) (*store.Store, error) { return s, nil }},
		&mockDeleteStore{nil},
	)
	resp, err := adapter.UpdateProfile(storeinterface.UpdateProfileParams{StoreID: 1, AddressID: &addr})
	if err != nil {
		t.Fatal(err)
	}
	if resp.AddressID == nil || *resp.AddressID != 99 {
		t.Error("unexpected address")
	}
}

func TestStoreAdapter_Delete_Success(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil},
		&mockDeleteStore{func(i appstore.DeleteStoreInput) error { return nil }},
	)
	err := adapter.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStoreAdapter_Delete_NotFound(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil}, &mockGetStore{nil}, &mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil},
		&mockDeleteStore{func(i appstore.DeleteStoreInput) error { return store.ErrStoreNotFound }},
	)
	err := adapter.Delete(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestStoreAdapter_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := storeinterface.NewAdapter(
		&mockCreateStore{nil},
		&mockGetStore{func(i appstore.GetStoreInput) (*store.Store, error) { return nil, unknownErr{} }},
		&mockListStores{nil},
		&mockToggleBulk{nil}, &mockToggleComm{nil}, &mockUpdateContact{nil},
		&mockUpdateName{nil}, &mockUpdateProfile{nil}, &mockDeleteStore{nil},
	)
	_, err := adapter.Get(1)
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
