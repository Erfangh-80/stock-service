package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	storeinterface "stock-service/internal/interface/store"
)

type StoreHandler struct {
	adapter *storeinterface.Adapter
}

func NewStoreHandler(adapter *storeinterface.Adapter) *StoreHandler {
	return &StoreHandler{adapter: adapter}
}

func (h *StoreHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/stores", h.Create)
	mux.HandleFunc("GET /api/v1/stores", h.List)
	mux.HandleFunc("GET /api/v1/stores/{id}", h.Get)
	mux.HandleFunc("PUT /api/v1/stores/{id}", h.UpdateName)
	mux.HandleFunc("PUT /api/v1/stores/{id}/contact", h.UpdateContact)
	mux.HandleFunc("PUT /api/v1/stores/{id}/profile", h.UpdateProfile)
	mux.HandleFunc("POST /api/v1/stores/{id}/bulk-sale", h.ToggleBulkSale)
	mux.HandleFunc("POST /api/v1/stores/{id}/commission", h.ToggleCommission)
	mux.HandleFunc("DELETE /api/v1/stores/{id}", h.Delete)
}

func (h *StoreHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(storeinterface.CreateStoreParams{
		UserID: body.UserID, StoreName: body.StoreName,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *StoreHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var userID *int64
	if v := q.Get("user_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
			return
		}
		userID = &id
	}

	var status *string
	if v := q.Get("status"); v != "" {
		status = &v
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.List(storeinterface.ListStoresFilter{
		UserID: userID,
		Status: status,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *StoreHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body dto.UpdateStoreNameRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.UpdateName(storeinterface.UpdateNameParams{
		StoreID: id, Name: body.Name,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body dto.UpdateStoreContactRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.UpdateContact(storeinterface.UpdateContactParams{
		StoreID: id, ContactPhone: body.ContactPhone,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body dto.UpdateStoreProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.UpdateProfile(storeinterface.UpdateProfileParams{
		StoreID:     id,
		AddressID:   body.AddressID,
		MediaAssets: body.MediaAssets,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreHandler) ToggleBulkSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.ToggleBulkSale(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreHandler) ToggleCommission(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.ToggleCommission(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
