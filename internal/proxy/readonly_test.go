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

// --- Unit tests for individual filter functions ---

func TestIsUpgrade(t *testing.T) {
	tests := []struct {
		name       string
		connection string
		upgrade    string
		want       bool
	}{
		{"no headers", "", "", false},
		{"Connection: Upgrade", "Upgrade", "", true},
		{"Upgrade: SPDY", "", "SPDY/3.1", true},
		{"Upgrade: websocket", "", "websocket", true},
		{"both headers", "Upgrade", "websocket", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.connection != "" {
				r.Header.Set("Connection", tt.connection)
			}
			if tt.upgrade != "" {
				r.Header.Set("Upgrade", tt.upgrade)
			}
			if got := isUpgrade(r); got != tt.want {
				t.Errorf("isUpgrade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsReadOnly(t *testing.T) {
	tests := []struct {
		method string
		want   bool
	}{
		{http.MethodGet, true},
		{http.MethodHead, true},
		{http.MethodOptions, true},
		{http.MethodPost, false},
		{http.MethodPut, false},
		{http.MethodPatch, false},
		{http.MethodDelete, false},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/v1/pods", nil)
			if got := isReadOnly(r); got != tt.want {
				t.Errorf("isReadOnly(%s) = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}

func TestIsNonMutatingPost(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   bool
	}{
		{"selfsubjectaccessreviews", http.MethodPost,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", true},
		{"subjectaccessreviews", http.MethodPost,
			"/apis/authorization.k8s.io/v1/subjectaccessreviews", true},
		{"localsubjectaccessreviews", http.MethodPost,
			"/apis/authorization.k8s.io/v1/namespaces/default/localsubjectaccessreviews", true},
		{"selfsubjectrulesreviews", http.MethodPost,
			"/apis/authorization.k8s.io/v1/selfsubjectrulesreviews", true},
		{"tokenreviews", http.MethodPost,
			"/apis/authentication.k8s.io/v1/tokenreviews", true},
		{"selfsubjectreviews", http.MethodPost,
			"/apis/authentication.k8s.io/v1/selfsubjectreviews", true},
		{"regular POST", http.MethodPost,
			"/api/v1/namespaces", false},
		{"GET to review path", http.MethodGet,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", false},
		{"DELETE to review path", http.MethodDelete,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", false},
		{"spoofed resource name", http.MethodPost,
			"/apis/evil.io/v1/selfsubjectaccessreviews", false},
		{"spoofed suffix in custom group", http.MethodPost,
			"/apis/custom.example.com/v1/namespaces/default/selfsubjectaccessreviews", false},
		{"review name as subresource", http.MethodPost,
			"/api/v1/namespaces/default/pods/selfsubjectaccessreviews", false},
		{"v1beta1 version allowed", http.MethodPost,
			"/apis/authorization.k8s.io/v1beta1/selfsubjectaccessreviews", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, nil)
			if got := isNonMutatingPost(r); got != tt.want {
				t.Errorf("isNonMutatingPost(%s %s) = %v, want %v", tt.method, tt.path, got, tt.want)
			}
		})
	}
}

func TestIsDryRun(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"dryRun=All", "/api/v1/namespaces?dryRun=All", true},
		{"no dryRun", "/api/v1/namespaces", false},
		{"dryRun=None", "/api/v1/namespaces?dryRun=None", false},
		{"dryRun empty", "/api/v1/namespaces?dryRun=", false},
		{"dryRun with other params", "/api/v1/namespaces?fieldManager=kubectl&dryRun=All", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			if got := isDryRun(r); got != tt.want {
				t.Errorf("isDryRun(%s) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

// --- Unit tests for the commander ---

func TestCheckRequest(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		headers map[string]string
		allowed bool
	}{
		{"GET allowed", http.MethodGet, "/api/v1/pods", nil, true},
		{"POST blocked", http.MethodPost, "/api/v1/pods", nil, false},
		{"upgrade blocked", http.MethodGet, "/api/v1/pods/foo/exec",
			map[string]string{"Connection": "Upgrade", "Upgrade": "SPDY/3.1"}, false},
		{"review POST allowed", http.MethodPost,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", nil, true},
		{"dry-run POST allowed", http.MethodPost,
			"/api/v1/namespaces?dryRun=All", nil, true},
		{"dry-run DELETE allowed", http.MethodDelete,
			"/api/v1/namespaces/foo?dryRun=All", nil, true},
		{"upgrade trumps dry-run", http.MethodGet, "/api/v1/pods?dryRun=All",
			map[string]string{"Connection": "Upgrade"}, false},
		{"upgrade trumps review", http.MethodPost,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews",
			map[string]string{"Connection": "Upgrade"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, nil)
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			reason, ok := checkRequest(r)
			if ok != tt.allowed {
				t.Errorf("checkRequest() allowed=%v, want %v (reason=%q)", ok, tt.allowed, reason)
			}
			if !ok && reason == "" {
				t.Error("checkRequest() returned blocked with empty reason")
			}
		})
	}
}

// --- Integration tests through the full handler ---

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

func TestHandler_AllowsNonMutatingPOST(t *testing.T) {
	handler, _ := newTestHandler(t)

	paths := []string{
		"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews",
		"/apis/authorization.k8s.io/v1/subjectaccessreviews",
		"/apis/authorization.k8s.io/v1/namespaces/default/localsubjectaccessreviews",
		"/apis/authorization.k8s.io/v1/selfsubjectrulesreviews",
		"/apis/authentication.k8s.io/v1/tokenreviews",
		"/apis/authentication.k8s.io/v1/selfsubjectreviews",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected 200 for POST %s, got %d", path, rr.Code)
			}
		})
	}
}

func TestHandler_AllowsDryRun(t *testing.T) {
	handler, _ := newTestHandler(t)

	for _, method := range []string{
		http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
	} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/namespaces?dryRun=All", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected 200 for %s with dryRun=All, got %d", method, rr.Code)
			}
		})
	}
}

func TestHandler_BlocksDryRunNone(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/namespaces?dryRun=None", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for dryRun=None, got %d", rr.Code)
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
