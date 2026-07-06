package warehouse_test

import (
	"errors"
	"testing"
	"time"

	"stock-service/internal/domain/warehouse"
	iface "stock-service/internal/interface"
	adapter "stock-service/internal/interface/warehouse"
)

// mocks

type mockCreateWarehouseUseCase struct {
	execute func(adapter.CreateWarehouseInput) (*warehouse.Warehouse, error)
}

func (m *mockCreateWarehouseUseCase) Execute(input adapter.CreateWarehouseInput) (*warehouse.Warehouse, error) {
	return m.execute(input)
}

type mockUpdateVisibilityUseCase struct {
	execute func(adapter.UpdateVisibilityInput) (*warehouse.Warehouse, error)
}

func (m *mockUpdateVisibilityUseCase) Execute(input adapter.UpdateVisibilityInput) (*warehouse.Warehouse, error) {
	return m.execute(input)
}

type mockUpdateContactUseCase struct {
	execute func(adapter.UpdateContactInput) (*warehouse.Warehouse, error)
}

func (m *mockUpdateContactUseCase) Execute(input adapter.UpdateContactInput) (*warehouse.Warehouse, error) {
	return m.execute(input)
}

// helpers

func newAdapter(
	create func(adapter.CreateWarehouseInput) (*warehouse.Warehouse, error),
	updateVis func(adapter.UpdateVisibilityInput) (*warehouse.Warehouse, error),
	updateCont func(adapter.UpdateContactInput) (*warehouse.Warehouse, error),
) *adapter.Adapter {
	return adapter.NewAdapter(
		&mockCreateWarehouseUseCase{execute: create},
		&mockUpdateVisibilityUseCase{execute: updateVis},
		&mockUpdateContactUseCase{execute: updateCont},
	)
}

func baseWarehouse() *warehouse.Warehouse {
	return &warehouse.Warehouse{
		ID:               1,
		CreatedByUserID:  100,
		WarehouseName:    "Main Warehouse",
		IsPublic:         true,
		CollectionMethod: "pickup",
		CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

// Test Create — success

func TestCreate_Success(t *testing.T) {
	a := newAdapter(
		func(input adapter.CreateWarehouseInput) (*warehouse.Warehouse, error) {
			if input.CreatedByUserID != 100 || input.WarehouseName != "Test Warehouse" {
				t.Errorf("unexpected input: %+v", input)
			}
			return baseWarehouse(), nil
		},
		nil, nil,
	)

	resp, err := a.Create(adapter.CreateWarehouseInput{CreatedByUserID: 100, WarehouseName: "Test Warehouse"})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 || resp.CreatedByUserID != 100 || resp.WarehouseName != "Main Warehouse" || !resp.IsPublic || resp.CollectionMethod != "pickup" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

// Test Create — ErrWarehouseNameRequired → ErrInvalidInput

func TestCreate_ErrWarehouseNameRequired(t *testing.T) {
	a := newAdapter(
		func(input adapter.CreateWarehouseInput) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNameRequired
		},
		nil, nil,
	)

	_, err := a.Create(adapter.CreateWarehouseInput{})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// Test Create — ErrWarehouseNameTooLong → ErrInvalidInput

func TestCreate_ErrWarehouseNameTooLong(t *testing.T) {
	a := newAdapter(
		func(input adapter.CreateWarehouseInput) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNameTooLong
		},
		nil, nil,
	)

	_, err := a.Create(adapter.CreateWarehouseInput{})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// Test Create — unknown error → ErrInternal

func TestCreate_UnknownError(t *testing.T) {
	a := newAdapter(
		func(input adapter.CreateWarehouseInput) (*warehouse.Warehouse, error) {
			return nil, errors.New("db connection failed")
		},
		nil, nil,
	)

	_, err := a.Create(adapter.CreateWarehouseInput{})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

// Test UpdateVisibility — success

func TestUpdateVisibility_Success(t *testing.T) {
	a := newAdapter(
		nil,
		func(input adapter.UpdateVisibilityInput) (*warehouse.Warehouse, error) {
			if input.WarehouseID != 1 || !input.IsPublic {
				t.Errorf("unexpected input: %+v", input)
			}
			return baseWarehouse(), nil
		},
		nil,
	)

	resp, err := a.UpdateVisibility(adapter.UpdateVisibilityInput{WarehouseID: 1, IsPublic: true})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

// Test UpdateVisibility — domain error → ErrInvalidInput

func TestUpdateVisibility_DomainError(t *testing.T) {
	a := newAdapter(
		nil,
		func(input adapter.UpdateVisibilityInput) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNameRequired
		},
		nil,
	)

	_, err := a.UpdateVisibility(adapter.UpdateVisibilityInput{WarehouseID: 1, IsPublic: true})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// Test UpdateVisibility — unknown error → ErrInternal

func TestUpdateVisibility_UnknownError(t *testing.T) {
	a := newAdapter(
		nil,
		func(input adapter.UpdateVisibilityInput) (*warehouse.Warehouse, error) {
			return nil, errors.New("unexpected")
		},
		nil,
	)

	_, err := a.UpdateVisibility(adapter.UpdateVisibilityInput{WarehouseID: 1, IsPublic: true})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

// Test UpdateContact — success

func TestUpdateContact_Success(t *testing.T) {
	phone := "123456"
	contactPhone := "789012"
	a := newAdapter(
		nil, nil,
		func(input adapter.UpdateContactInput) (*warehouse.Warehouse, error) {
			if input.WarehouseID != 1 || *input.Phone != "123456" || *input.ContactPhone != "789012" || input.CollectionMethod != "delivery" {
				t.Errorf("unexpected input: %+v", input)
			}
			return baseWarehouse(), nil
		},
	)

	resp, err := a.UpdateContact(adapter.UpdateContactInput{
		WarehouseID: 1, Phone: &phone, ContactPhone: &contactPhone, CollectionMethod: "delivery",
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

// Test UpdateContact — domain error → ErrInvalidInput

func TestUpdateContact_DomainError(t *testing.T) {
	a := newAdapter(
		nil, nil,
		func(input adapter.UpdateContactInput) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNameTooLong
		},
	)

	_, err := a.UpdateContact(adapter.UpdateContactInput{WarehouseID: 1})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// Test UpdateContact — unknown error → ErrInternal

func TestUpdateContact_UnknownError(t *testing.T) {
	a := newAdapter(
		nil, nil,
		func(input adapter.UpdateContactInput) (*warehouse.Warehouse, error) {
			return nil, errors.New("disk full")
		},
	)

	_, err := a.UpdateContact(adapter.UpdateContactInput{WarehouseID: 1})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
