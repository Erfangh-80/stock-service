package storewarehouselink_test

import (
	"errors"
	"testing"

	"stock-service/internal/domain/store_warehouse_link"
	iface "stock-service/internal/interface"
	adapter "stock-service/internal/interface/store_warehouse_link"
)

// mocks

type mockCreateLinkUseCase struct {
	execute func(adapter.CreateLinkInput) (*storewarehouselink.StoreWarehouseLink, error)
}

func (m *mockCreateLinkUseCase) Execute(input adapter.CreateLinkInput) (*storewarehouselink.StoreWarehouseLink, error) {
	return m.execute(input)
}

type mockChangeRelationUseCase struct {
	execute func(adapter.ChangeRelationInput) (*storewarehouselink.StoreWarehouseLink, error)
}

func (m *mockChangeRelationUseCase) Execute(input adapter.ChangeRelationInput) (*storewarehouselink.StoreWarehouseLink, error) {
	return m.execute(input)
}

// helpers

func newAdapter(
	create func(adapter.CreateLinkInput) (*storewarehouselink.StoreWarehouseLink, error),
	change func(adapter.ChangeRelationInput) (*storewarehouselink.StoreWarehouseLink, error),
) *adapter.Adapter {
	return adapter.NewAdapter(
		&mockCreateLinkUseCase{execute: create},
		&mockChangeRelationUseCase{execute: change},
	)
}

func ptr[T any](v T) *T { return &v }

// Test Create — success

func TestCreate_Success(t *testing.T) {
	a := newAdapter(
		func(input adapter.CreateLinkInput) (*storewarehouselink.StoreWarehouseLink, error) {
			if input.StoreID != 10 || input.WarehouseID != 20 {
				t.Errorf("unexpected input: %+v", input)
			}
			return &storewarehouselink.StoreWarehouseLink{
				ID: 1, StoreID: 10, WarehouseID: 20, RelationType: storewarehouselink.RelationTypePrimary,
			}, nil
		},
		nil,
	)

	resp, err := a.Create(adapter.CreateLinkInput{StoreID: 10, WarehouseID: 20})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 || resp.StoreID != 10 || resp.WarehouseID != 20 || resp.RelationType != "primary" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

// Test Create — error

func TestCreate_Error(t *testing.T) {
	someErr := errors.New("oops")
	a := newAdapter(
		func(input adapter.CreateLinkInput) (*storewarehouselink.StoreWarehouseLink, error) {
			return nil, someErr
		},
		nil,
	)

	_, err := a.Create(adapter.CreateLinkInput{StoreID: 1, WarehouseID: 2})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

// Test ChangeRelation — success

func TestChangeRelation_Success(t *testing.T) {
	a := newAdapter(
		nil,
		func(input adapter.ChangeRelationInput) (*storewarehouselink.StoreWarehouseLink, error) {
			if input.LinkID != 5 || input.RelationType != "primary" {
				t.Errorf("unexpected input: %+v", input)
			}
			return &storewarehouselink.StoreWarehouseLink{
				ID: 5, StoreID: 10, WarehouseID: 20, RelationType: storewarehouselink.RelationTypePrimary,
			}, nil
		},
	)

	resp, err := a.ChangeRelation(adapter.ChangeRelationInput{LinkID: 5, RelationType: "primary"})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 5 || resp.StoreID != 10 || resp.WarehouseID != 20 || resp.RelationType != "primary" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

// Test ChangeRelation — error

func TestChangeRelation_Error(t *testing.T) {
	someErr := errors.New("oops")
	a := newAdapter(
		nil,
		func(input adapter.ChangeRelationInput) (*storewarehouselink.StoreWarehouseLink, error) {
			return nil, someErr
		},
	)

	_, err := a.ChangeRelation(adapter.ChangeRelationInput{LinkID: 1, RelationType: "primary"})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

// Test unknown error → ErrInternal (redundant with above but makes intent explicit)

func TestCreate_UnknownErrorMapsToInternal(t *testing.T) {
	a := newAdapter(
		func(input adapter.CreateLinkInput) (*storewarehouselink.StoreWarehouseLink, error) {
			return nil, errors.New("some random error")
		},
		nil,
	)

	_, err := a.Create(adapter.CreateLinkInput{StoreID: 1, WarehouseID: 2})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestChangeRelation_UnknownErrorMapsToInternal(t *testing.T) {
	a := newAdapter(
		nil,
		func(input adapter.ChangeRelationInput) (*storewarehouselink.StoreWarehouseLink, error) {
			return nil, errors.New("some random error")
		},
	)

	_, err := a.ChangeRelation(adapter.ChangeRelationInput{LinkID: 1, RelationType: "primary"})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
