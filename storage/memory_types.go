package storage

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	mu      sync.Mutex
	entries map[string]entry
}
