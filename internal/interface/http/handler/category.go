package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	categoryinterface "stock-service/internal/interface/category"
)

type CategoryHandler struct {
	adapter *categoryinterface.Adapter
}

func NewCategoryHandler(adapter *categoryinterface.Adapter) *CategoryHandler {
	return &CategoryHandler{adapter: adapter}
}

func (h *CategoryHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/categories", h.Create)
	mux.HandleFunc("GET /api/v1/categories", h.List)
	mux.HandleFunc("GET /api/v1/categories/{id}", h.Get)
	mux.HandleFunc("PUT /api/v1/categories/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/categories/{id}", h.Delete)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(categoryinterface.CreateCategoryParams{
		Name: body.Name, Slug: body.Slug,
		ParentID: body.ParentID, Description: body.Description,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.adapter.List()
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Update(categoryinterface.UpdateCategoryParams{
		ID: id, Name: body.Name, Slug: body.Slug,
		Description: body.Description, ParentID: body.ParentID,
		SortOrder: body.SortOrder,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
