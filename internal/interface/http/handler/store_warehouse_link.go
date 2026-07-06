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
	mux.HandleFunc("GET /api/v1/warehouse-links", h.List)
	mux.HandleFunc("GET /api/v1/warehouse-links/{id}", h.Get)
	mux.HandleFunc("DELETE /api/v1/warehouse-links/{id}", h.Delete)
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

func (h *StoreWarehouseLinkHandler) List(w http.ResponseWriter, r *http.Request) {
	var storeID *int64
	if sid := r.URL.Query().Get("store_id"); sid != "" {
		if v, err := strconv.ParseInt(sid, 10, 64); err == nil {
			storeID = &v
		}
	}
	var warehouseID *int64
	if wid := r.URL.Query().Get("warehouse_id"); wid != "" {
		if v, err := strconv.ParseInt(wid, 10, 64); err == nil {
			warehouseID = &v
		}
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	result, err := h.adapter.List(storewarehouselink.ListLinksInput{
		StoreID: storeID, WarehouseID: warehouseID, Page: page, Limit: limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreWarehouseLinkHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.Get(storewarehouselink.GetLinkInput{ID: id})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreWarehouseLinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.adapter.Delete(storewarehouselink.DeleteLinkInput{ID: id}); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
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
