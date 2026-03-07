package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// setupTestServer creates a lightweight HTTP server with stub endpoints for testing.
func setupTestServer(t *testing.T) (*httptest.Server, *Client) {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-1"}}}`))
	})

	mux.HandleFunc("/users/uid-123/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"currentLevel":5,"currentState":{"type":"on"}}`))
			return
		}
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Retry") == "done" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}
		w.WriteHeader(http.StatusTooManyRequests)
	})

	srv := httptest.NewServer(mux)

	// Pre-populate a valid token so auth is skipped during tests
	relaxtechClient := New("email", "pass", "", "", "")
	relaxtechClient.BaseURL = srv.URL
	relaxtechClient.token = "t"
	relaxtechClient.tokenExp = time.Now().Add(time.Hour)
	relaxtechClient.HTTP = srv.Client()

	return srv, relaxtechClient
}

func TestUserIDAutoPopulated(t *testing.T) {
	srv, relaxtechClient := setupTestServer(t)
	defer srv.Close()

	// UserID starts empty; GetStatus should resolve it via /users/me
	st, err := relaxtechClient.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if relaxtechClient.UserID != "uid-123" {
		t.Fatalf("expected user id populated, got %s", relaxtechClient.UserID)
	}
	if st.CurrentLevel != 5 || st.CurrentState.Type != "on" {
		t.Fatalf("unexpected status %+v", st)
	}
}

func TestRateLimitRetry(t *testing.T) {
	requestCount := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	relaxtechClient := New("email", "pass", "uid", "", "")
	relaxtechClient.BaseURL = srv.URL
	relaxtechClient.token = "t"
	relaxtechClient.tokenExp = time.Now().Add(time.Hour)
	relaxtechClient.HTTP = srv.Client()

	start := time.Now()
	if err := relaxtechClient.do(context.Background(), http.MethodGet, "/ping", nil, nil, nil); err != nil {
		t.Fatalf("do retry: %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("expected 2 attempts, got %d", requestCount)
	}
	if elapsed := time.Since(start); elapsed < 2*time.Second {
		t.Fatalf("expected backoff delay, got %v", elapsed)
	}
}
