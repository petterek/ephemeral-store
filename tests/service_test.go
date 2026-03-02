package tests

import (
	"testing"
	"time"

	"ephemral-storage/storage"
)

func TestInsertAndReadOnce(t *testing.T) {
	svc := storage.NewService(storage.NewMemoryStore())

	if err := svc.InsertKeyValue("k1", "secret", 60); err != nil {
		t.Fatalf("InsertKeyValue: %v", err)
	}

	val, err := svc.ReadValue("k1")
	if err != nil {
		t.Fatalf("ReadValue: %v", err)
	}
	if val != "secret" {
		t.Fatalf("got %q, want %q", val, "secret")
	}

	// Second read must fail
	_, err = svc.ReadValue("k1")
	if err != storage.ErrKeyNotFound {
		t.Fatalf("second ReadValue: got %v, want ErrKeyNotFound", err)
	}
}

func TestReadExpired(t *testing.T) {
	svc := storage.NewService(storage.NewMemoryStore())

	if err := svc.InsertKeyValue("k2", "gone", 1); err != nil {
		t.Fatalf("InsertKeyValue: %v", err)
	}

	time.Sleep(1100 * time.Millisecond)

	_, err := svc.ReadValue("k2")
	if err != storage.ErrKeyNotFound {
		t.Fatalf("expired ReadValue: got %v, want ErrKeyNotFound", err)
	}
}

func TestReadMissing(t *testing.T) {
	svc := storage.NewService(storage.NewMemoryStore())

	_, err := svc.ReadValue("nope")
	if err != storage.ErrKeyNotFound {
		t.Fatalf("missing ReadValue: got %v, want ErrKeyNotFound", err)
	}
}
