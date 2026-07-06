package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	brandinterface "stock-service/internal/interface/brand"
)

type BrandHandler struct {
	adapter *brandinterface.Adapter
}

func NewBrandHandler(adapter *brandinterface.Adapter) *BrandHandler {
	return &BrandHandler{adapter: adapter}
}

func (h *BrandHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/brands", h.Create)
	mux.HandleFunc("GET /api/v1/brands", h.List)
	mux.HandleFunc("GET /api/v1/brands/{id}", h.Get)
	mux.HandleFunc("PUT /api/v1/brands/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/brands/{id}", h.Delete)
}

func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(brandinterface.CreateBrandParams{
		Name: body.Name, Slug: body.Slug,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *BrandHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.adapter.List()
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *BrandHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.Get(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body dto.UpdateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Update(brandinterface.UpdateBrandParams{
		ID: id, Name: body.Name, Slug: body.Slug,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.adapter.Delete(id); err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
