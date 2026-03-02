package storage

// Store defines the interface for ephemeral key-value storage backends.
type Store interface {
	// Set stores a value with a TTL in seconds.
	Set(key, value string, ttlSeconds int) error
	// GetAndDelete retrieves a value and removes it so it can only be read once.
	// Returns the value and true if found, or "" and false if not.
	GetAndDelete(key string) (string, bool)
	// List returns all non-expired entries as key-value pairs.
	List() []KeyValue
}

// KeyValue represents a stored entry with its remaining TTL.
type KeyValue struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ExpiresIn int    `json:"expires_in"`
}
