package proxy

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTestHandler(t *testing.T) (http.Handler, *httptest.Server) {
	t.Helper()
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-Method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	t.Cleanup(backend.Close)

	target, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(target, http.DefaultTransport)
	return handler, backend
}

func TestHandler_AllowedMethods(t *testing.T) {
	handler, _ := newTestHandler(t)

	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/pods", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected 200, got %d", rr.Code)
			}
		})
	}
}

func TestHandler_BlockedMethods(t *testing.T) {
	handler, _ := newTestHandler(t)

	for _, method := range []string{
		http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch,
	} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/pods", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rr.Code)
			}

			var status metav1.Status
			if err := json.NewDecoder(rr.Body).Decode(&status); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if status.Status != metav1.StatusFailure {
				t.Errorf("expected status Failure, got %q", status.Status)
			}
			if status.Reason != metav1.StatusReasonMethodNotAllowed {
				t.Errorf("expected reason MethodNotAllowed, got %q", status.Reason)
			}
			if status.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected code 405, got %d", status.Code)
			}
		})
	}
}

func TestHandler_BlocksUpgrade(t *testing.T) {
	handler, _ := newTestHandler(t)

	tests := []struct {
		name       string
		connection string
		upgrade    string
	}{
		{"SPDY upgrade", "Upgrade", "SPDY/3.1"},
		{"WebSocket upgrade", "Upgrade", "websocket"},
		{"Upgrade header only", "", "websocket"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/pods/foo/exec", nil)
			if tt.connection != "" {
				req.Header.Set("Connection", tt.connection)
			}
			if tt.upgrade != "" {
				req.Header.Set("Upgrade", tt.upgrade)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rr.Code)
			}
		})
	}
}

func TestHandler_GETResponsePassthrough(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"kind":"PodList","items":[]}`))
	}))
	t.Cleanup(backend.Close)

	target, _ := url.Parse(backend.URL)
	handler := NewHandler(target, http.DefaultTransport)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	body, _ := io.ReadAll(rr.Body)
	if string(body) != `{"kind":"PodList","items":[]}` {
		t.Errorf("unexpected response body: %s", body)
	}
}
