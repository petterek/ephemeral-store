package storage

import (
	"sync"
)

type entry struct {
	value     string
	expiresAt int64 // unix nanoseconds
}

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	mu      sync.RWMutex
	entries map[string]entry
}
