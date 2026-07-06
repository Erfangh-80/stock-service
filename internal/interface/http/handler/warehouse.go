package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	warehouseinterface "stock-service/internal/interface/warehouse"
)

type WarehouseHandler struct {
	adapter *warehouseinterface.Adapter
}

func NewWarehouseHandler(adapter *warehouseinterface.Adapter) *WarehouseHandler {
	return &WarehouseHandler{adapter: adapter}
}

func (h *WarehouseHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/warehouses", h.Create)
	mux.HandleFunc("PUT /api/v1/warehouses/{id}/visibility", h.UpdateVisibility)
	mux.HandleFunc("PUT /api/v1/warehouses/{id}/contact", h.UpdateContact)
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Create(warehouseinterface.CreateWarehouseInput{
		CreatedByUserID: body.CreatedByUserID,
		WarehouseName:   body.WarehouseName,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *WarehouseHandler) UpdateVisibility(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.UpdateVisibilityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.UpdateVisibility(warehouseinterface.UpdateVisibilityInput{
		WarehouseID: id, IsPublic: body.IsPublic,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *WarehouseHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.UpdateWarehouseContactRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.UpdateContact(warehouseinterface.UpdateContactInput{
		WarehouseID: id, Phone: body.Phone,
		ContactPhone: body.ContactPhone, CollectionMethod: body.CollectionMethod,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}
