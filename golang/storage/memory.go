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
	exp := time.Now().UnixNano() + int64(ttlSeconds)*1e9
	m.mu.Lock()
	m.entries[key] = entry{
		value:     value,
		expiresAt: exp,
	}
	m.mu.Unlock()
	return nil
}

func (m *MemoryStore) GetAndDelete(key string) (string, bool) {
	m.mu.Lock()
	e, ok := m.entries[key]
	if !ok {
		m.mu.Unlock()
		return "", false
	}
	delete(m.entries, key)
	m.mu.Unlock()
	if time.Now().UnixNano() > e.expiresAt {
		return "", false
	}
	return e.value, true
}

func (m *MemoryStore) List() []KeyValue {
	m.mu.Lock()
	now := time.Now().UnixNano()
	result := make([]KeyValue, 0, len(m.entries))
	for k, e := range m.entries {
		if now > e.expiresAt {
			delete(m.entries, k)
			continue
		}
		result = append(result, KeyValue{
			Key:       k,
			Value:     e.value,
			ExpiresIn: int((e.expiresAt - now) / 1e9),
		})
	}
	m.mu.Unlock()
	return result
}
