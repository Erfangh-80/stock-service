package warehouse_test

import (
	"errors"
	"testing"
	"time"

	app "stock-service/internal/application/warehouse"
	"stock-service/internal/domain/warehouse"
	iface "stock-service/internal/interface"
	adapter "stock-service/internal/interface/warehouse"
)

type mockCreateWH struct {
	fn func(int64, string) (*warehouse.Warehouse, error)
}

func (m *mockCreateWH) Execute(createdByUserID int64, warehouseName string) (*warehouse.Warehouse, error) {
	return m.fn(createdByUserID, warehouseName)
}

type mockGetWH struct {
	fn func(int64) (*warehouse.Warehouse, error)
}

func (m *mockGetWH) Execute(warehouseID int64) (*warehouse.Warehouse, error) {
	return m.fn(warehouseID)
}

type mockListWH struct {
	fn func(app.ListWarehousesInput) (*app.ListWarehousesOutput, error)
}

func (m *mockListWH) Execute(input app.ListWarehousesInput) (*app.ListWarehousesOutput, error) {
	if m.fn != nil {
		return m.fn(input)
	}
	return &app.ListWarehousesOutput{}, nil
}

type mockDelWH struct{}

func (m *mockDelWH) Execute(warehouseID int64) error { return nil }

type mockVisWH struct{}

func (m *mockVisWH) Execute(warehouseID int64, isPublic bool) error { return nil }

type mockContWH struct{}

func (m *mockContWH) Execute(warehouseID int64, phone, contactPhone *string, collectionMethod string) error {
	return nil
}

type mockUpdWH struct {
	fn func(app.UpdateWarehouseInput) (*warehouse.Warehouse, error)
}

func (m *mockUpdWH) Execute(input app.UpdateWarehouseInput) (*warehouse.Warehouse, error) {
	if m.fn != nil {
		return m.fn(input)
	}
	return baseWarehouse(), nil
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

func newAdapter(
	create *mockCreateWH,
	get *mockGetWH,
	list *mockListWH,
	del *mockDelWH,
	vis *mockVisWH,
	cont *mockContWH,
	upd *mockUpdWH,
) *adapter.Adapter {
	return adapter.NewAdapter(create, get, list, del, vis, cont, upd)
}

func TestCreate_Success(t *testing.T) {
	a := newAdapter(
		&mockCreateWH{fn: func(uid int64, name string) (*warehouse.Warehouse, error) {
			if uid != 100 || name != "Test Warehouse" {
				t.Errorf("unexpected input: uid=%d name=%s", uid, name)
			}
			return baseWarehouse(), nil
		}},
		nil, nil, nil, nil, nil, nil,
	)
	resp, err := a.Create(adapter.CreateWarehouseInput{CreatedByUserID: 100, WarehouseName: "Test Warehouse"})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 || resp.CreatedByUserID != 100 || resp.WarehouseName != "Main Warehouse" || !resp.IsPublic || resp.CollectionMethod != "pickup" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestCreate_ErrWarehouseNameRequired(t *testing.T) {
	a := newAdapter(
		&mockCreateWH{fn: func(uid int64, name string) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNameRequired
		}},
		nil, nil, nil, nil, nil, nil,
	)
	_, err := a.Create(adapter.CreateWarehouseInput{})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreate_ErrWarehouseNameTooLong(t *testing.T) {
	a := newAdapter(
		&mockCreateWH{fn: func(uid int64, name string) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNameTooLong
		}},
		nil, nil, nil, nil, nil, nil,
	)
	_, err := a.Create(adapter.CreateWarehouseInput{})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreate_UnknownError(t *testing.T) {
	a := newAdapter(
		&mockCreateWH{fn: func(uid int64, name string) (*warehouse.Warehouse, error) {
			return nil, errors.New("db connection failed")
		}},
		nil, nil, nil, nil, nil, nil,
	)
	_, err := a.Create(adapter.CreateWarehouseInput{})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestGet_Success(t *testing.T) {
	a := newAdapter(
		nil,
		&mockGetWH{fn: func(id int64) (*warehouse.Warehouse, error) {
			if id != 1 {
				t.Errorf("unexpected id: %d", id)
			}
			return baseWarehouse(), nil
		}},
		nil, nil, nil, nil, nil,
	)
	resp, err := a.Get(adapter.GetWarehouseInput{ID: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGet_NotFound(t *testing.T) {
	a := newAdapter(
		nil,
		&mockGetWH{fn: func(id int64) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNotFound
		}},
		nil, nil, nil, nil, nil,
	)
	_, err := a.Get(adapter.GetWarehouseInput{ID: 999})
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGet_UnknownError(t *testing.T) {
	a := newAdapter(
		nil,
		&mockGetWH{fn: func(id int64) (*warehouse.Warehouse, error) {
			return nil, errors.New("unexpected")
		}},
		nil, nil, nil, nil, nil,
	)
	_, err := a.Get(adapter.GetWarehouseInput{ID: 1})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	a := newAdapter(nil, nil, nil, &mockDelWH{}, nil, nil, nil)
	err := a.Delete(adapter.DeleteWarehouseInput{ID: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestUpdateVisibility_Success(t *testing.T) {
	a := newAdapter(
		nil,
		&mockGetWH{fn: func(id int64) (*warehouse.Warehouse, error) { return baseWarehouse(), nil }},
		nil, nil, &mockVisWH{}, nil, nil,
	)
	resp, err := a.UpdateVisibility(adapter.UpdateVisibilityInput{WarehouseID: 1, IsPublic: true})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

func TestUpdateContact_Success(t *testing.T) {
	a := newAdapter(
		nil,
		&mockGetWH{fn: func(id int64) (*warehouse.Warehouse, error) { return baseWarehouse(), nil }},
		nil, nil, nil, &mockContWH{}, nil,
	)
	resp, err := a.UpdateContact(adapter.UpdateContactInput{WarehouseID: 1, CollectionMethod: "pickup"})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

func TestUpdateWarehouse_Success(t *testing.T) {
	a := newAdapter(
		nil, nil, nil, nil, nil, nil,
		&mockUpdWH{fn: func(input app.UpdateWarehouseInput) (*warehouse.Warehouse, error) {
			if input.WarehouseID != 1 {
				t.Errorf("unexpected id: %d", input.WarehouseID)
			}
			return baseWarehouse(), nil
		}},
	)
	name := "Updated"
	resp, err := a.UpdateWarehouse(adapter.UpdateWarehouseInput{WarehouseID: 1, Name: &name})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

func TestUpdateWarehouse_ErrWarehouseNameRequired(t *testing.T) {
	a := newAdapter(
		nil, nil, nil, nil, nil, nil,
		&mockUpdWH{fn: func(input app.UpdateWarehouseInput) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNameRequired
		}},
	)
	empty := ""
	_, err := a.UpdateWarehouse(adapter.UpdateWarehouseInput{WarehouseID: 1, Name: &empty})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestList_Success(t *testing.T) {
	a := newAdapter(
		nil, nil,
		&mockListWH{fn: func(input app.ListWarehousesInput) (*app.ListWarehousesOutput, error) {
			return &app.ListWarehousesOutput{
				Warehouses: []*warehouse.Warehouse{baseWarehouse()},
				Total:      1,
			}, nil
		}},
		nil, nil, nil, nil,
	)
	resp, err := a.List(adapter.ListWarehousesInput{Page: 1, Limit: 20})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(resp.Warehouses) != 1 || resp.Total != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestUpdateWarehouse_ErrWarehouseNotFound(t *testing.T) {
	a := newAdapter(
		nil, nil, nil, nil, nil, nil,
		&mockUpdWH{fn: func(input app.UpdateWarehouseInput) (*warehouse.Warehouse, error) {
			return nil, warehouse.ErrWarehouseNotFound
		}},
	)
	_, err := a.UpdateWarehouse(adapter.UpdateWarehouseInput{WarehouseID: 999})
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
