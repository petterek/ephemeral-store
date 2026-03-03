package tests

import (
	"fmt"
	"sync"
	"testing"

	"ephemral-storage/storage"
)

func BenchmarkInsert(b *testing.B) {
	svc := storage.NewService(storage.NewMemoryStore())
	for i := 0; i < b.N; i++ {
		svc.InsertKeyValue(fmt.Sprintf("key-%d", i), "value", 60)
	}
}

func BenchmarkInsertAndRead(b *testing.B) {
	svc := storage.NewService(storage.NewMemoryStore())
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		svc.InsertKeyValue(key, "value", 60)
		svc.ReadValue(key)
	}
}

func BenchmarkConcurrentInsert(b *testing.B) {
	svc := storage.NewService(storage.NewMemoryStore())
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			svc.InsertKeyValue(fmt.Sprintf("key-%d", i), "value", 60)
			i++
		}
	})
}

func BenchmarkConcurrentInsertAndRead(b *testing.B) {
	svc := storage.NewService(storage.NewMemoryStore())
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i)
			svc.InsertKeyValue(key, "value", 60)
			svc.ReadValue(key)
			i++
		}
	})
}

func BenchmarkList(b *testing.B) {
	store := storage.NewMemoryStore()
	svc := storage.NewService(store)
	for i := 0; i < 1000; i++ {
		svc.InsertKeyValue(fmt.Sprintf("key-%d", i), "value", 300)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.List()
	}
}

func BenchmarkHTTPInsertAndRead(b *testing.B) {
	svc := storage.NewService(storage.NewMemoryStore())
	hub := storage.NewHub(svc)

	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			svc.InsertKeyValue(key, "value", 60)
			hub.Broadcast()
			svc.ReadValue(key)
			hub.Broadcast()
		}(i)
	}
	wg.Wait()
}
