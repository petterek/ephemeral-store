package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Hub manages SSE clients and broadcasts storage state updates.
type Hub struct {
	mu      sync.Mutex
	clients map[chan []byte]struct{}
	svc     *Service
}

// NewHub creates a Hub tied to the given Service.
func NewHub(svc *Service) *Hub {
	return &Hub{
		clients: make(map[chan []byte]struct{}),
		svc:     svc,
	}
}

// Broadcast sends the current storage state to all connected SSE clients.
func (h *Hub) Broadcast() {
	entries := h.svc.store.List()
	if entries == nil {
		entries = []KeyValue{}
	}
	data, _ := json.Marshal(entries)

	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- data:
		default:
		}
	}
}

// HandleSSE serves an SSE stream of storage state.
func (h *Hub) HandleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan []byte, 8)

	// Send initial state
	entries := h.svc.store.List()
	if entries == nil {
		entries = []KeyValue{}
	}
	initial, _ := json.Marshal(entries)
	fmt.Fprintf(w, "data: %s\n\n", initial)
	flusher.Flush()

	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, ch)
		h.mu.Unlock()
	}()

	for {
		select {
		case data := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
