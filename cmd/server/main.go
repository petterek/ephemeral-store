package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"ephemral-storage/storage"
)

//go:embed static
var staticFiles embed.FS

func main() {
	store := storage.NewMemoryStore()
	svc := storage.NewService(store)
	hub := storage.NewHub(svc)
	handler := storage.NewHandler(svc, hub)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
