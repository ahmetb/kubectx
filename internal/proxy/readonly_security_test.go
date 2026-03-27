package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ============================================================
// Security & Jailbreak Test Suite for kubectx readonly proxy
// ============================================================

// --- Jailbreak: HTTP method smuggling ---

func TestJailbreak_MethodOverrideHeaders(t *testing.T) {
	handler, _ := newTestHandler(t)

	// Attackers might try X-HTTP-Method-Override or similar headers
	// to smuggle a POST through as a GET.
	overrideHeaders := []string{
		"X-HTTP-Method-Override",
		"X-HTTP-Method",
		"X-Method-Override",
	}

	for _, hdr := range overrideHeaders {
		t.Run(hdr, func(t *testing.T) {
			// Send GET with override header claiming POST
			req := httptest.NewRequest(http.MethodGet, "/api/v1/namespaces", nil)
			req.Header.Set(hdr, "POST")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// The proxy should allow it (it's a GET), but the backend should
			// NOT see it as a POST. The key is: do these headers reach the backend?
			// If backend honors the override header, the readonly protection is bypassed.
			if rr.Code == http.StatusOK {
				// Check if backend received the override header - it could be dangerous
				t.Logf("INFO: %s header passed through to backend (method override headers are forwarded)", hdr)
			}
		})
	}
}

// --- Jailbreak: Path traversal and URL manipulation ---

func TestJailbreak_PathTraversal(t *testing.T) {
	handler, _ := newTestHandler(t)

	paths := []struct {
		name   string
		method string
		path   string
		want   int
	}{
		// Try to sneak a POST through with path encoding
		{"encoded POST path", http.MethodPost, "/api/v1/%6eamespaces", 405},
		{"double-encoded path", http.MethodPost, "/api/v1/%256eamespaces", 405},
		// Try making a DELETE look like a review endpoint
		{"DELETE masquerading as review", http.MethodDelete,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", 405},
		// PUT masquerading as review
		{"PUT masquerading as review", http.MethodPut,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", 405},
		// PATCH masquerading as review
		{"PATCH masquerading as review", http.MethodPatch,
			"/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", 405},
	}

	for _, tt := range paths {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.want {
				t.Errorf("expected %d, got %d", tt.want, rr.Code)
			}
		})
	}
}

// --- Jailbreak: DryRun parameter smuggling ---

func TestJailbreak_DryRunSmuggling(t *testing.T) {
	handler, _ := newTestHandler(t)

	tests := []struct {
		name    string
		method  string
		rawURL  string
		wantOK  bool
		comment string
	}{
		{"dryRun=All (legit)", http.MethodPost, "/api/v1/namespaces?dryRun=All", true,
			"legitimate dry-run should work"},
		{"dryRun=all lowercase", http.MethodPost, "/api/v1/namespaces?dryRun=all", false,
			"case-sensitive: 'all' != 'All'"},
		{"dryRun=ALL uppercase", http.MethodPost, "/api/v1/namespaces?dryRun=ALL", false,
			"case-sensitive: 'ALL' != 'All'"},
		{"dryRun=All+extra", http.MethodPost, "/api/v1/namespaces?dryRun=All&dryRun=None", true,
			"POTENTIAL ISSUE: multiple dryRun params - Go Query().Get() returns first"},
		{"dryRun=None+All", http.MethodPost, "/api/v1/namespaces?dryRun=None&dryRun=All", false,
			"Query().Get() returns first param which is None - should be blocked"},
		{"dryRun with spaces", http.MethodPost, "/api/v1/namespaces?dryRun=%20All", false,
			"space-padded dryRun should be rejected"},
		{"dryRun=All with null byte", http.MethodPost, "/api/v1/namespaces?dryRun=All%00None", false,
			"null byte injection in dryRun value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.rawURL, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			gotOK := rr.Code == http.StatusOK
			if gotOK != tt.wantOK {
				t.Errorf("%s: expected ok=%v, got status %d [%s]", tt.name, tt.wantOK, rr.Code, tt.comment)
			} else {
				t.Logf("PASS: %s [%s]", tt.name, tt.comment)
			}
		})
	}
}

// --- Jailbreak: Upgrade header smuggling ---

func TestJailbreak_UpgradeSmuggling(t *testing.T) {
	handler, _ := newTestHandler(t)

	tests := []struct {
		name    string
		method  string
		path    string
		headers map[string]string
		wantOK  bool
	}{
		// Case variation - Go canonicalizes header keys, so "connection" → "Connection"
		{"connection key lowercase (Go canonicalizes)", http.MethodGet, "/api/v1/pods/x/exec",
			map[string]string{"connection": "Upgrade"}, false,
			// Go normalizes header key to "Connection", value "Upgrade" matches exactly → blocked
		},
		// Fixed: Connection value "upgrade" (lowercase) now caught by strings.EqualFold
		{"Connection: upgrade value lowercase", http.MethodGet, "/api/v1/pods/x/exec",
			map[string]string{"Connection": "upgrade"}, false,
		},
		// Multiple Connection header values
		{"Connection with multiple values", http.MethodGet, "/api/v1/pods/x/exec",
			map[string]string{"Connection": "keep-alive, Upgrade", "Upgrade": "SPDY/3.1"}, false,
			// Has Upgrade header set so isUpgrade catches it via second check
		},
		// Empty upgrade header
		{"empty Upgrade header", http.MethodGet, "/api/v1/pods/x/exec",
			map[string]string{"Upgrade": ""}, true,
			// Empty string != "" is false, so this passes through
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			gotOK := rr.Code == http.StatusOK
			if gotOK != tt.wantOK {
				t.Errorf("expected ok=%v, got status %d", tt.wantOK, rr.Code)
			}
		})
	}
}

// --- Jailbreak: Review endpoint spoofing ---

func TestJailbreak_ReviewEndpointSpoofing(t *testing.T) {
	handler, _ := newTestHandler(t)

	tests := []struct {
		name   string
		path   string
		wantOK bool
	}{
		// Legitimate
		{"legit SSAR", "/apis/authorization.k8s.io/v1/selfsubjectaccessreviews", true},
		// Trailing slash
		{"trailing slash", "/apis/authorization.k8s.io/v1/selfsubjectaccessreviews/", false},
		// Path with query string
		{"with query string", "/apis/authorization.k8s.io/v1/selfsubjectaccessreviews?foo=bar", true},
		// CRD in different API group matching same resource name
		{"spoofed API group", "/apis/authorization.evil.io/v1/selfsubjectaccessreviews", false},
		// Subresource under a review endpoint name
		{"subresource", "/apis/authorization.k8s.io/v1/selfsubjectaccessreviews/status", false},
		// Double dot in API group
		{"double-dot group", "/apis/authorization..k8s..io/v1/selfsubjectaccessreviews", false},
		// Unicode homograph
		{"unicode homograph", "/apis/authorization.k8s.іo/v1/selfsubjectaccessreviews", false},
		// Namespace-scoped version of cluster-scoped review
		{"namespace-scoped SSAR", "/apis/authorization.k8s.io/v1/namespaces/default/selfsubjectaccessreviews", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			gotOK := rr.Code == http.StatusOK
			if gotOK != tt.wantOK {
				t.Errorf("expected ok=%v, got status %d", tt.wantOK, rr.Code)
			}
		})
	}
}

// --- Jailbreak: CONNECT method (HTTP tunneling) ---

func TestJailbreak_ConnectMethod(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest("CONNECT", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("CONNECT should be blocked, got %d", rr.Code)
	}
}

// --- Jailbreak: Custom/unusual HTTP methods ---

func TestJailbreak_UnusualMethods(t *testing.T) {
	handler, _ := newTestHandler(t)

	methods := []string{
		"CONNECT",
		"TRACE",
		"PROPFIND",     // WebDAV
		"MKCOL",        // WebDAV
		"COPY",         // WebDAV
		"MOVE",         // WebDAV
		"LOCK",         // WebDAV
		"UNLOCK",       // WebDAV
		"PURGE",        // Varnish
		"LINK",         // Link
		"UNLINK",       // Unlink
		"VIEW",         // non-standard
		"CUSTOMDELETE", // non-standard
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/pods", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("%s should be blocked, got %d", method, rr.Code)
			}
		})
	}
}

// --- Jailbreak: Large request body on allowed endpoint ---

func TestJailbreak_LargeBodyOnGET(t *testing.T) {
	handler, _ := newTestHandler(t)

	// Some proxies might convert a GET with a body to a POST
	body := strings.NewReader(`{"kind":"Namespace","apiVersion":"v1","metadata":{"name":"evil"}}`)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/namespaces", body)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// GET with body should still be treated as GET (allowed)
	if rr.Code != http.StatusOK {
		t.Errorf("GET with body should still be allowed, got %d", rr.Code)
	}
}

// --- Blocked response format verification ---

func TestBlockedResponse_Format(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pods", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Verify Content-Type
	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	// Verify response body is valid Kubernetes Status
	var status metav1.Status
	body, _ := io.ReadAll(rr.Body)
	if err := json.Unmarshal(body, &status); err != nil {
		t.Fatalf("response is not valid JSON: %v\nbody: %s", err, body)
	}

	if status.APIVersion != "v1" {
		t.Errorf("expected apiVersion=v1, got %q", status.APIVersion)
	}
	if status.Kind != "Status" {
		t.Errorf("expected kind=Status, got %q", status.Kind)
	}
	if status.Code != 405 {
		t.Errorf("expected code=405, got %d", status.Code)
	}
	if status.Message == "" {
		t.Error("expected non-empty message")
	}
	if !strings.Contains(status.Message, "[kubectx]") {
		t.Errorf("expected message to contain [kubectx], got %q", status.Message)
	}
}

// --- Kubeconfig rewriting security tests ---

func TestRewriteKubeconfig_Security(t *testing.T) {
	// Verify that no credentials leak into the rewritten kubeconfig
	input := `
apiVersion: v1
kind: Config
clusters:
- name: test-cluster
  cluster:
    server: https://real-server.example.com:6443
    certificate-authority-data: c2VjcmV0LWNh
contexts:
- name: test-ctx
  context:
    cluster: test-cluster
    user: test-user
current-context: test-ctx
users:
- name: test-user
  user:
    client-certificate-data: c2VjcmV0LWNlcnQ=
    client-key-data: c2VjcmV0LWtleQ==
    token: super-secret-token
`
	out, err := RewriteKubeconfig([]byte(input), "127.0.0.1:12345")
	if err != nil {
		t.Fatal(err)
	}

	result := string(out)

	// Real server URL must not appear
	if strings.Contains(result, "real-server.example.com") {
		t.Error("SECURITY: real server URL leaked into rewritten kubeconfig")
	}

	// Credentials must not appear
	sensitiveStrings := []string{
		"c2VjcmV0LWNh",       // CA data
		"c2VjcmV0LWNlcnQ=",   // client cert
		"c2VjcmV0LWtleQ==",   // client key
		"super-secret-token",  // token
	}
	for _, s := range sensitiveStrings {
		if strings.Contains(result, s) {
			t.Errorf("SECURITY: credential data %q leaked into rewritten kubeconfig", s)
		}
	}

	// Proxy address must be present
	if !strings.Contains(result, "127.0.0.1:12345") {
		t.Error("proxy address not found in rewritten kubeconfig")
	}

	// [RO] suffix must be present
	if !strings.Contains(result, "[RO]") {
		t.Error("[RO] context suffix not found in rewritten kubeconfig")
	}
}

// --- Concurrent request handling ---

func TestConcurrent_MixedRequests(t *testing.T) {
	handler, _ := newTestHandler(t)

	done := make(chan struct{}, 100)

	// Fire off mixed GET and POST requests concurrently
	for i := 0; i < 50; i++ {
		go func(i int) {
			defer func() { done <- struct{}{} }()
			var method string
			if i%2 == 0 {
				method = http.MethodGet
			} else {
				method = http.MethodPost
			}
			req := httptest.NewRequest(method, "/api/v1/pods", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if method == http.MethodGet && rr.Code != http.StatusOK {
				t.Errorf("GET #%d: expected 200, got %d", i, rr.Code)
			}
			if method == http.MethodPost && rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("POST #%d: expected 405, got %d", i, rr.Code)
			}
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}
}

// --- Watchlist (GET with watch param - should be allowed) ---

func TestWatch_AllowedViaGET(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pods?watch=true", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET with watch=true should be allowed, got %d", rr.Code)
	}
}

// --- kubectl apply --dry-run=server sends POST with dryRun=All ---

func TestDryRunApply_ServerSide(t *testing.T) {
	handler, _ := newTestHandler(t)

	// Simulates: kubectl apply --dry-run=server
	body := strings.NewReader(`{"kind":"Namespace","apiVersion":"v1","metadata":{"name":"test"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/namespaces?dryRun=All&fieldManager=kubectl-client-side-apply", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("server-side dry-run apply should be allowed, got %d", rr.Code)
	}
}

// --- kubectl logs should work (GET, no upgrade) ---

func TestLogs_AllowedViaGET(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/namespaces/default/pods/my-pod/log?follow=true", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("kubectl logs should be allowed (GET), got %d", rr.Code)
	}
}

// --- kubectl exec should be blocked ---

func TestExec_Blocked(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/namespaces/default/pods/my-pod/exec?command=sh", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "SPDY/3.1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("kubectl exec should be blocked, got %d", rr.Code)
	}

	var status metav1.Status
	json.NewDecoder(rr.Body).Decode(&status)
	if !strings.Contains(status.Message, "exec") {
		t.Errorf("blocked message should mention exec, got %q", status.Message)
	}
}

// --- kubectl cp should be blocked ---

func TestCp_Blocked(t *testing.T) {
	handler, _ := newTestHandler(t)

	// kubectl cp first does exec
	req := httptest.NewRequest(http.MethodPost, "/api/v1/namespaces/default/pods/my-pod/exec?command=tar", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "SPDY/3.1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("kubectl cp should be blocked, got %d", rr.Code)
	}
}

// --- kubectl port-forward should be blocked ---

func TestPortForward_Blocked(t *testing.T) {
	handler, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/namespaces/default/pods/my-pod/portforward", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "SPDY/3.1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("kubectl port-forward should be blocked, got %d", rr.Code)
	}
}

// --- E2E proxy integration test with real HTTP server ---

func TestE2E_ProxyWithBackend(t *testing.T) {
	var backendRequests []string

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backendRequests = append(backendRequests, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch {
		case r.URL.Path == "/api/v1/pods" && r.Method == "GET":
			w.Write([]byte(`{"kind":"PodList","apiVersion":"v1","items":[{"metadata":{"name":"test-pod"}}]}`))
		case r.URL.Path == "/api/v1/namespaces" && r.Method == "GET":
			w.Write([]byte(`{"kind":"NamespaceList","apiVersion":"v1","items":[{"metadata":{"name":"default"}}]}`))
		default:
			w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Success"}`))
		}
	}))
	defer backend.Close()

	target, _ := url.Parse(backend.URL)
	handler := NewHandler(target, http.DefaultTransport)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	client := proxyServer.Client()

	// Test 1: GET pods - should reach backend
	resp, err := client.Get(proxyServer.URL + "/api/v1/pods")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("GET pods: expected 200, got %d", resp.StatusCode)
	}
	if !strings.Contains(string(body), "test-pod") {
		t.Error("GET pods: response doesn't contain expected pod")
	}

	// Test 2: POST namespace - should be blocked (never reach backend)
	preCount := len(backendRequests)
	resp, err = client.Post(proxyServer.URL+"/api/v1/namespaces", "application/json",
		strings.NewReader(`{"kind":"Namespace","apiVersion":"v1","metadata":{"name":"evil"}}`))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Errorf("POST namespace: expected 405, got %d", resp.StatusCode)
	}
	if len(backendRequests) != preCount {
		t.Error("POST namespace: request leaked through to backend!")
	}

	// Test 3: DELETE pod - should be blocked
	req, _ := http.NewRequest(http.MethodDelete, proxyServer.URL+"/api/v1/namespaces/default/pods/test-pod", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 405 {
		t.Errorf("DELETE pod: expected 405, got %d", resp.StatusCode)
	}

	// Test 4: GET namespaces - should work
	resp, err = client.Get(proxyServer.URL + "/api/v1/namespaces")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("GET namespaces: expected 200, got %d", resp.StatusCode)
	}

	// Test 5: dry-run POST - should reach backend
	req, _ = http.NewRequest(http.MethodPost, proxyServer.URL+"/api/v1/namespaces?dryRun=All", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("dry-run POST: expected 200, got %d", resp.StatusCode)
	}
}
