package application

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	service *FeatureFlagService
}

func NewHandler(service *FeatureFlagService) *Handler {
	return &Handler{service: service}
}

type CreateFlagRequest struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	GlobalEnabled bool   `json:"globalEnabled"`
}

type SetOverrideRequest struct {
	Enabled bool `json:"enabled"`
}

type SetGlobalRequest struct {
	Enabled bool `json:"enabled"`
}

type EvaluateResponse struct {
	Key     string `json:"key"`
	UserID  string `json:"userId"`
	Enabled bool   `json:"enabled"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	var req CreateFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Key == "" || req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "key and name are required"})
		return
	}

	flag, err := h.service.CreateFlag(req.Key, req.Name, req.GlobalEnabled)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flag)
}

func (h *Handler) GetFlag(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	flag, err := h.service.GetFlag(key)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flag)
}

func (h *Handler) SetUserOverride(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	userID := chi.URLParam(r, "userId")

	var req SetOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	flag, err := h.service.SetUserOverride(key, userID, req.Enabled)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flag)
}

func (h *Handler) SetGlobalState(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	var req SetGlobalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	flag, err := h.service.SetGlobalState(key, req.Enabled)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flag)
}

func (h *Handler) EvaluateFlag(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	userID := r.URL.Query().Get("userId")

	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "userId query parameter is required"})
		return
	}

	enabled, err := h.service.Evaluate(key, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(EvaluateResponse{
		Key:     key,
		UserID:  userID,
		Enabled: enabled,
	})
}

func (h *Handler) ListFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := h.service.ListFlags()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flags)
}
