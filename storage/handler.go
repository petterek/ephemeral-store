package storage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Handler provides HTTP handlers for the ephemeral storage service.
type Handler struct {
	svc *Service
	hub *Hub
}

// NewHandler creates a Handler backed by the given Service and Hub.
func NewHandler(svc *Service, hub *Hub) *Handler {
	return &Handler{svc: svc, hub: hub}
}

// RegisterRoutes registers the API routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /keys", h.handleInsert)
	mux.HandleFunc("GET /keys/{key}", h.handleRead)
	mux.HandleFunc("GET /events", h.hub.HandleSSE)
	mux.HandleFunc("GET /isalive", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /isready", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
}

func (h *Handler) handleInsert(w http.ResponseWriter, r *http.Request) {
	var req insertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	if req.Sender == "" || req.DataType == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "sender and datatype are required"})
		return
	}
	if req.TTL <= 0 {
		req.TTL = 60
	}
	raw := fmt.Sprintf("%s.%s.%d", req.Sender, req.DataType, time.Now().UnixMilli())
	key := base64.URLEncoding.EncodeToString([]byte(raw))
	if err := h.svc.InsertKeyValue(key, req.Value, req.TTL); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, insertResponse{Key: key})
	h.hub.Broadcast()
}

func (h *Handler) handleRead(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	val, err := h.svc.ReadValue(key)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, readResponse{Value: val})
	h.hub.Broadcast()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
