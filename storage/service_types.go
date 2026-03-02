package storage

import "errors"

// ErrKeyNotFound is returned when a key does not exist or has already been read.
var ErrKeyNotFound = errors.New("key not found")

// Service provides ephemeral key-value operations backed by a Store.
type Service struct {
	store Store
}
