package api

import (
	"bookstore/recommendation/internal/store"
	"bookstore/recommendation/internal/workflow"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestHTTPRegisterAndHealth(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	server := New(workflow.New(st), nil)
	request := httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(`{"id":"api-1","store_id":"s","title":"API Book","score":77}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("register status: %d", response.Code)
	}
	health := httptest.NewRecorder()
	server.Handler().ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health status: %d", health.Code)
	}
}
