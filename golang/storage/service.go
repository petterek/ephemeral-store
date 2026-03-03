package storage

// NewService creates a Service with the given Store backend.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// InsertKeyValue stores a value that expires after ttlSeconds.
func (s *Service) InsertKeyValue(key, value string, ttl int) error {
	return s.store.Set(key, value, ttl)
}

// ReadValue retrieves a value exactly once. Subsequent reads return an error.
func (s *Service) ReadValue(key string) (string, error) {
	val, ok := s.store.GetAndDelete(key)
	if !ok {
		return "", ErrKeyNotFound
	}
	return val, nil
}
