package storewarehouselink

import (
	"stock-service/internal/domain/store_warehouse_link"
	iface "stock-service/internal/interface"
)



type CreateLinkInput struct {
	StoreID     int64
	WarehouseID int64
}

type ChangeRelationInput struct {
	LinkID       int64
	RelationType string
}

type LinkResponse struct {
	ID           int64  `json:"id"`
	StoreID      int64  `json:"store_id"`
	WarehouseID  int64  `json:"warehouse_id"`
	RelationType string `json:"relation_type"`
}

type CreateLinkUseCase interface {
	Execute(input CreateLinkInput) (*storewarehouselink.StoreWarehouseLink, error)
}

type ChangeRelationUseCase interface {
	Execute(input ChangeRelationInput) (*storewarehouselink.StoreWarehouseLink, error)
}

type Adapter struct {
	create  CreateLinkUseCase
	change  ChangeRelationUseCase
}

func NewAdapter(create CreateLinkUseCase, change ChangeRelationUseCase) *Adapter {
	return &Adapter{create: create, change: change}
}

func (a *Adapter) Create(input CreateLinkInput) (*LinkResponse, error) {
	swl, err := a.create.Execute(input)
	if err != nil {
		return nil, iface.ErrInternal
	}
	return &LinkResponse{
		ID:           swl.ID,
		StoreID:      swl.StoreID,
		WarehouseID:  swl.WarehouseID,
		RelationType: string(swl.RelationType),
	}, nil
}

func (a *Adapter) ChangeRelation(input ChangeRelationInput) (*LinkResponse, error) {
	swl, err := a.change.Execute(input)
	if err != nil {
		return nil, iface.ErrInternal
	}
	return &LinkResponse{
		ID:           swl.ID,
		StoreID:      swl.StoreID,
		WarehouseID:  swl.WarehouseID,
		RelationType: string(swl.RelationType),
	}, nil
}
