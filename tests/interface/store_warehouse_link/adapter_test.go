package storewarehouselink_test

import (
	"testing"

	app "stock-service/internal/application/store_warehouse_link"
	domain "stock-service/internal/domain/store_warehouse_link"
	iface "stock-service/internal/interface"
	adapter "stock-service/internal/interface/store_warehouse_link"
)

type mockCreateLinkUseCase struct {
	fn func(int64, int64) (*domain.StoreWarehouseLink, error)
}

func (m *mockCreateLinkUseCase) Execute(storeID, warehouseID int64) (*domain.StoreWarehouseLink, error) {
	return m.fn(storeID, warehouseID)
}

type mockGetLinkUseCase struct {
	fn func(app.GetLinkInput) (*domain.StoreWarehouseLink, error)
}

func (m *mockGetLinkUseCase) Execute(input app.GetLinkInput) (*domain.StoreWarehouseLink, error) {
	return m.fn(input)
}

type mockListLinksUseCase struct {
	fn func(app.ListLinksInput) (*app.ListLinksOutput, error)
}

func (m *mockListLinksUseCase) Execute(input app.ListLinksInput) (*app.ListLinksOutput, error) {
	return m.fn(input)
}

type mockChangeRelationUseCase struct {
	fn func(app.ChangeRelationInput) (*domain.StoreWarehouseLink, error)
}

func (m *mockChangeRelationUseCase) Execute(input app.ChangeRelationInput) (*domain.StoreWarehouseLink, error) {
	return m.fn(input)
}

type mockDeleteLinkUseCase struct {
	fn func(app.DeleteLinkInput) error
}

func (m *mockDeleteLinkUseCase) Execute(input app.DeleteLinkInput) error {
	return m.fn(input)
}

func newTestAdapter(
	create *mockCreateLinkUseCase,
	get *mockGetLinkUseCase,
	list *mockListLinksUseCase,
	change *mockChangeRelationUseCase,
	del *mockDeleteLinkUseCase,
) *adapter.Adapter {
	if create == nil {
		create = &mockCreateLinkUseCase{}
	}
	if get == nil {
		get = &mockGetLinkUseCase{}
	}
	if list == nil {
		list = &mockListLinksUseCase{}
	}
	if change == nil {
		change = &mockChangeRelationUseCase{}
	}
	if del == nil {
		del = &mockDeleteLinkUseCase{}
	}
	return adapter.NewAdapter(create, get, list, change, del)
}

func TestCreate_Success(t *testing.T) {
	a := newTestAdapter(
		&mockCreateLinkUseCase{func(storeID, warehouseID int64) (*domain.StoreWarehouseLink, error) {
			return &domain.StoreWarehouseLink{
				ID: 1, StoreID: storeID, WarehouseID: warehouseID,
				RelationType: domain.RelationTypePrimary,
			}, nil
		}},
		nil, nil, nil, nil,
	)

	resp, err := a.Create(adapter.CreateLinkInput{StoreID: 10, WarehouseID: 20})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 || resp.StoreID != 10 || resp.WarehouseID != 20 || resp.RelationType != "primary" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestCreate_Error(t *testing.T) {
	a := newTestAdapter(
		&mockCreateLinkUseCase{func(storeID, warehouseID int64) (*domain.StoreWarehouseLink, error) {
			return nil, iface.ErrInternal
		}},
		nil, nil, nil, nil,
	)

	_, err := a.Create(adapter.CreateLinkInput{StoreID: 1, WarehouseID: 2})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestGet_Success(t *testing.T) {
	a := newTestAdapter(
		nil,
		&mockGetLinkUseCase{func(input app.GetLinkInput) (*domain.StoreWarehouseLink, error) {
			return &domain.StoreWarehouseLink{
				ID: input.ID, StoreID: 10, WarehouseID: 20, RelationType: domain.RelationTypePrimary,
			}, nil
		}},
		nil, nil, nil,
	)

	resp, err := a.Get(adapter.GetLinkInput{ID: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestGet_NotFound_ReturnsNotFound(t *testing.T) {
	a := newTestAdapter(
		nil,
		&mockGetLinkUseCase{func(input app.GetLinkInput) (*domain.StoreWarehouseLink, error) {
			return nil, domain.ErrLinkNotFound
		}},
		nil, nil, nil,
	)

	_, err := a.Get(adapter.GetLinkInput{ID: 999})
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList_Success(t *testing.T) {
	a := newTestAdapter(
		nil, nil,
		&mockListLinksUseCase{func(input app.ListLinksInput) (*app.ListLinksOutput, error) {
			return &app.ListLinksOutput{
				Links: []*domain.StoreWarehouseLink{
					{ID: 1, StoreID: 10, WarehouseID: 20, RelationType: domain.RelationTypePrimary},
				},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil,
	)

	resp, err := a.List(adapter.ListLinksInput{})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(resp.Links) != 1 || resp.Total != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestChangeRelation_Success(t *testing.T) {
	a := newTestAdapter(
		nil, nil, nil,
		&mockChangeRelationUseCase{func(input app.ChangeRelationInput) (*domain.StoreWarehouseLink, error) {
			return &domain.StoreWarehouseLink{
				ID: input.LinkID, StoreID: 10, WarehouseID: 20,
				RelationType: domain.RelationTypePrimary,
			}, nil
		}},
		nil,
	)

	resp, err := a.ChangeRelation(adapter.ChangeRelationInput{LinkID: 5, RelationType: "primary"})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if resp.ID != 5 || resp.RelationType != "primary" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestChangeRelation_InvalidType_ReturnsInvalidInput(t *testing.T) {
	a := newTestAdapter(
		nil, nil, nil,
		&mockChangeRelationUseCase{func(input app.ChangeRelationInput) (*domain.StoreWarehouseLink, error) {
			return nil, domain.ErrInvalidRelationType
		}},
		nil,
	)

	_, err := a.ChangeRelation(adapter.ChangeRelationInput{LinkID: 1, RelationType: "invalid"})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	a := newTestAdapter(
		nil, nil, nil, nil,
		&mockDeleteLinkUseCase{func(input app.DeleteLinkInput) error {
			return nil
		}},
	)

	err := a.Delete(adapter.DeleteLinkInput{ID: 1})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestDelete_UnknownError_ReturnsInternal(t *testing.T) {
	a := newTestAdapter(
		nil, nil, nil, nil,
		&mockDeleteLinkUseCase{func(input app.DeleteLinkInput) error {
			return iface.ErrInternal
		}},
	)

	err := a.Delete(adapter.DeleteLinkInput{ID: 1})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
