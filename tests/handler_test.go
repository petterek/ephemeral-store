package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ephemral-storage/storage"
)

func setupTestServer() *httptest.Server {
	svc := storage.NewService(storage.NewMemoryStore())
	hub := storage.NewHub(svc)
	handler := storage.NewHandler(svc, hub)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

func TestHTTPInsertAndRead(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// Insert
	resp, err := http.Post(ts.URL+"/keys", "application/json",
		strings.NewReader(`{"sender":"alice","datatype":"password","value":"secret","ttl":60}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("insert: got %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	var insertResp struct{ Key string `json:"key"` }
	json.NewDecoder(resp.Body).Decode(&insertResp)
	resp.Body.Close()
	if insertResp.Key == "" {
		t.Fatal("insert: expected key in response")
	}

	// First read succeeds
	resp, err = http.Get(ts.URL + "/keys/" + insertResp.Key)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("first read: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	resp.Body.Close()

	// Second read returns 404
	resp, err = http.Get(ts.URL + "/keys/" + insertResp.Key)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("second read: got %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
	resp.Body.Close()
}

func TestHTTPReadMissing(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/keys/nope")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("missing: got %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
	resp.Body.Close()
}

func TestHTTPInsertBadRequest(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/keys", "application/json",
		strings.NewReader(`{"sender":"","datatype":"","value":"v"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad request: got %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
	resp.Body.Close()
}
