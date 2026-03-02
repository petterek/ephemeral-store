package storage

import (
	"time"
)

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		entries: make(map[string]entry),
	}
}

func (m *MemoryStore) Set(key, value string, ttlSeconds int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
	return nil
}

func (m *MemoryStore) GetAndDelete(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[key]
	if !ok {
		return "", false
	}
	delete(m.entries, key)
	if time.Now().After(e.expiresAt) {
		return "", false
	}
	return e.value, true
}

func (m *MemoryStore) List() []KeyValue {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	var result []KeyValue
	for k, e := range m.entries {
		if now.After(e.expiresAt) {
			delete(m.entries, k)
			continue
		}
		result = append(result, KeyValue{
			Key:       k,
			Value:     e.value,
			ExpiresIn: int(time.Until(e.expiresAt).Seconds()),
		})
	}
	return result
}
