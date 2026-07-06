package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	storewarehouselink "stock-service/internal/interface/store_warehouse_link"
)

type StoreWarehouseLinkHandler struct {
	adapter *storewarehouselink.Adapter
}

func NewStoreWarehouseLinkHandler(adapter *storewarehouselink.Adapter) *StoreWarehouseLinkHandler {
	return &StoreWarehouseLinkHandler{adapter: adapter}
}

func (h *StoreWarehouseLinkHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/warehouse-links", h.Create)
	mux.HandleFunc("PUT /api/v1/warehouse-links/{id}/relation", h.ChangeRelation)
}

func (h *StoreWarehouseLinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateWarehouseLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Create(storewarehouselink.CreateLinkInput{
		StoreID: body.StoreID, WarehouseID: body.WarehouseID,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *StoreWarehouseLinkHandler) ChangeRelation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.ChangeRelationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.ChangeRelation(storewarehouselink.ChangeRelationInput{
		LinkID: id, RelationType: body.RelationType,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}
