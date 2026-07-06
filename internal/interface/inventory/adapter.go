package inventoryinterface

import (
	"time"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
	"stock-service/internal/domain/product"
	iface "stock-service/internal/interface"
)

type InventoryResponse struct {
	ID               int64
	StoreID          int64
	WarehouseID      int64
	ProductID        int32
	SaleModel        string
	BasePrice        float64
	PromotionID      *int64
	FinalPrice       *float64
	StartAt          *time.Time
	EndAt            *time.Time
	PromotionStatus  string
	InstantQty       int
	MinOrderQty      int
	MaxOrderQty      *int
	Condition        string
	VendorSaleStatus string
	SystemSaleStatus string
	CreatedAt        time.Time
}

type InventoryListResponse struct {
	Items []InventoryResponse `json:"items"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

type LowStockResponse struct {
	IsLow      bool `json:"is_low"`
	CurrentQty int  `json:"current_qty"`
}

type CreateInventoryParams struct {
	StoreID     int64
	WarehouseID int64
	ProductID   int32
	BasePrice   float64
}

type ApplyPromotionParams struct {
	SaleID      int64
	PromotionID int64
	FinalPrice  float64
	StartAt     time.Time
	EndAt       time.Time
}

type RemovePromotionParams struct {
	SaleID int64
}

type UpdateInventoryParams struct {
	SaleID       int64
	InstantQty   int
	ScheduledQty map[string]int
	MinOrderQty  int
	MaxOrderQty  *int
}

type ListInventoryParams struct {
	StoreID          *int64
	ProductID        *int32
	VendorSaleStatus *string
	SystemSaleStatus *string
	Page             int
	Limit            int
}

type DeleteInventoryParams struct {
	SaleID int64
}

type SearchInventoryParams struct {
	Query string
	Page  int
	Limit int
}

type ReserveQuantityParams struct {
	SaleID   int64
	Quantity int
}

type ReleaseQuantityParams struct {
	SaleID   int64
	Quantity int
}

type CheckLowStockParams struct {
	SaleID    int64
	Threshold int
}

type Adapter struct {
	create         appinventory.CreateInventoryUseCase
	get            appinventory.GetInventoryUseCase
	list           appinventory.ListInventoryUseCase
	delete         appinventory.DeleteInventoryUseCase
	search         appinventory.SearchInventoryUseCase
	apply          appinventory.ApplyPromotionUseCase
	remove         appinventory.RemovePromotionUseCase
	update         appinventory.UpdateInventoryUseCase
	suspendVendor  appinventory.SuspendVendorSaleUseCase
	closeVendor    appinventory.CloseVendorSaleUseCase
	suspendSystem  appinventory.SuspendSystemSaleUseCase
	closeSystem    appinventory.CloseSystemSaleUseCase
	reserveQty     appinventory.ReserveQuantityUseCase
	releaseQty     appinventory.ReleaseQuantityUseCase
	checkLowStock  appinventory.CheckLowStockUseCase
}

func NewAdapter(
	create appinventory.CreateInventoryUseCase,
	get appinventory.GetInventoryUseCase,
	list appinventory.ListInventoryUseCase,
	delete appinventory.DeleteInventoryUseCase,
	search appinventory.SearchInventoryUseCase,
	apply appinventory.ApplyPromotionUseCase,
	remove appinventory.RemovePromotionUseCase,
	update appinventory.UpdateInventoryUseCase,
	suspendVendor appinventory.SuspendVendorSaleUseCase,
	closeVendor appinventory.CloseVendorSaleUseCase,
	suspendSystem appinventory.SuspendSystemSaleUseCase,
	closeSystem appinventory.CloseSystemSaleUseCase,
	reserveQty appinventory.ReserveQuantityUseCase,
	releaseQty appinventory.ReleaseQuantityUseCase,
	checkLowStock appinventory.CheckLowStockUseCase,
) *Adapter {
	return &Adapter{
		create: create, get: get, list: list, delete: delete, search: search,
		apply: apply, remove: remove, update: update,
		suspendVendor: suspendVendor, closeVendor: closeVendor,
		suspendSystem: suspendSystem, closeSystem: closeSystem,
		reserveQty: reserveQty, releaseQty: releaseQty,
		checkLowStock: checkLowStock,
	}
}

func (a *Adapter) Create(params CreateInventoryParams) (*InventoryResponse, error) {
	result, err := a.create.Execute(appinventory.CreateInventoryInput{
		StoreID: params.StoreID, WarehouseID: params.WarehouseID,
		ProductID: params.ProductID, BasePrice: params.BasePrice,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Get(id int64) (*InventoryResponse, error) {
	result, err := a.get.Execute(appinventory.GetInventoryInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) List(params ListInventoryParams) (*InventoryListResponse, error) {
	result, err := a.list.Execute(appinventory.ListInventoryInput{
		StoreID: params.StoreID, ProductID: params.ProductID,
		VendorSaleStatus: params.VendorSaleStatus, SystemSaleStatus: params.SystemSaleStatus,
		Page: params.Page, Limit: params.Limit,
	})
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]InventoryResponse, len(result.Items))
	for i, inv := range result.Items {
		items[i] = *toResponse(inv)
	}
	return &InventoryListResponse{Items: items, Total: result.Total, Page: result.Page, Limit: result.Limit}, nil
}

func (a *Adapter) Delete(saleID int64) error {
	err := a.delete.Execute(appinventory.DeleteInventoryInput{SaleID: saleID})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Search(params SearchInventoryParams) (*InventoryListResponse, error) {
	result, err := a.search.Execute(appinventory.SearchInventoryInput{
		Query: params.Query, Page: params.Page, Limit: params.Limit,
	})
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]InventoryResponse, len(result.Items))
	for i, inv := range result.Items {
		items[i] = *toResponse(inv)
	}
	return &InventoryListResponse{Items: items, Total: result.Total, Page: result.Page, Limit: result.Limit}, nil
}

func (a *Adapter) ApplyPromotion(params ApplyPromotionParams) (*InventoryResponse, error) {
	result, err := a.apply.Execute(appinventory.ApplyPromotionInput{
		SaleID: params.SaleID, PromotionID: params.PromotionID,
		FinalPrice: params.FinalPrice, StartAt: params.StartAt, EndAt: params.EndAt,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) RemovePromotion(params RemovePromotionParams) (*InventoryResponse, error) {
	result, err := a.remove.Execute(appinventory.RemovePromotionInput{SaleID: params.SaleID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) UpdateInventory(params UpdateInventoryParams) (*InventoryResponse, error) {
	result, err := a.update.Execute(appinventory.UpdateInventoryInput{
		SaleID: params.SaleID, InstantQty: params.InstantQty,
		ScheduledQty: params.ScheduledQty, MinOrderQty: params.MinOrderQty,
		MaxOrderQty: params.MaxOrderQty,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) SuspendVendorSale(saleID int64) (*InventoryResponse, error) {
	result, err := a.suspendVendor.Execute(appinventory.SuspendVendorSaleInput{SaleID: saleID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) CloseVendorSale(saleID int64) (*InventoryResponse, error) {
	result, err := a.closeVendor.Execute(appinventory.CloseVendorSaleInput{SaleID: saleID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) SuspendSystemSale(saleID int64) (*InventoryResponse, error) {
	result, err := a.suspendSystem.Execute(appinventory.SuspendSystemSaleInput{SaleID: saleID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) CloseSystemSale(saleID int64) (*InventoryResponse, error) {
	result, err := a.closeSystem.Execute(appinventory.CloseSystemSaleInput{SaleID: saleID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) ReserveQuantity(saleID int64, qty int) (*InventoryResponse, error) {
	result, err := a.reserveQty.Execute(appinventory.ReserveQuantityInput{SaleID: saleID, Quantity: qty})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) ReleaseQuantity(saleID int64, qty int) (*InventoryResponse, error) {
	result, err := a.releaseQty.Execute(appinventory.ReleaseQuantityInput{SaleID: saleID, Quantity: qty})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) CheckLowStock(saleID int64, threshold int) (*LowStockResponse, error) {
	result, err := a.checkLowStock.Execute(appinventory.CheckLowStockInput{SaleID: saleID, Threshold: threshold})
	if err != nil {
		return nil, mapError(err)
	}
	return &LowStockResponse{IsLow: result.IsLow, CurrentQty: result.CurrentQty}, nil
}

func mapError(err error) error {
	switch err {
	case inventory.ErrInventoryNotFound, product.ErrProductNotFound:
		return iface.ErrNotFound
	case inventory.ErrInvalidBasePrice, inventory.ErrInvalidFinalPrice,
		inventory.ErrInvalidQuantity, inventory.ErrInvalidMinOrderQty,
		inventory.ErrInvalidMaxOrderQty, inventory.ErrInvalidPromotionDates,
		inventory.ErrInvalidScheduledDate,
		inventory.ErrVendorSaleStatusTransition, inventory.ErrSystemSaleStatusTransition:
		return iface.ErrInvalidInput
	case inventory.ErrPromotionAlreadyApplied, inventory.ErrNoActivePromotion:
		return iface.ErrConflict
	case inventory.ErrInsufficientStock:
		return iface.ErrConflict
	default:
		return iface.ErrInternal
	}
}

func toResponse(inv *inventory.Inventory) *InventoryResponse {
	var maxOrderQty *int
	if inv.MaxOrderQty != nil {
		v := *inv.MaxOrderQty
		maxOrderQty = &v
	}

	return &InventoryResponse{
		ID: inv.ID, StoreID: inv.StoreID, WarehouseID: inv.WarehouseID,
		ProductID: inv.ProductID, SaleModel: string(inv.SaleModel),
		BasePrice: inv.BasePrice, PromotionID: inv.PromotionID,
		FinalPrice: inv.FinalPrice, StartAt: inv.StartAt, EndAt: inv.EndAt,
		PromotionStatus: string(inv.PromotionStatus), InstantQty: inv.InstantQty,
		MinOrderQty: inv.MinOrderQty, MaxOrderQty: maxOrderQty,
		Condition: string(inv.Condition), VendorSaleStatus: string(inv.VendorSaleStatus),
		SystemSaleStatus: string(inv.SystemSaleStatus), CreatedAt: inv.CreatedAt,
	}
}
