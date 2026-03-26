package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ahmetb/kubectx/internal/env"
)

// nonMutatingPostPatterns match Kubernetes "review" endpoints that use POST
// but don't create persistent resources. Patterns are anchored to known API
// groups to prevent spoofing via custom resource names.
var nonMutatingPostPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^/apis/authorization\.k8s\.io/[^/]+/selfsubjectaccessreviews$`),
	regexp.MustCompile(`^/apis/authorization\.k8s\.io/[^/]+/subjectaccessreviews$`),
	regexp.MustCompile(`^/apis/authorization\.k8s\.io/[^/]+/namespaces/[^/]+/localsubjectaccessreviews$`),
	regexp.MustCompile(`^/apis/authorization\.k8s\.io/[^/]+/selfsubjectrulesreviews$`),
	regexp.MustCompile(`^/apis/authentication\.k8s\.io/[^/]+/tokenreviews$`),
	regexp.MustCompile(`^/apis/authentication\.k8s\.io/[^/]+/selfsubjectreviews$`),
}

var debugLog = func() *log.Logger {
	if _, ok := os.LookupEnv(env.EnvDebug); ok {
		return log.New(os.Stderr, "[readonly-proxy] ", log.Ltime)
	}
	return log.New(nopWriter{}, "", 0)
}()

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

// ReadonlyProxy is a reverse proxy that only allows read-only HTTP methods.
type ReadonlyProxy struct {
	server   *http.Server
	listener net.Listener
}

// Config holds information needed to start the readonly proxy.
type Config struct {
	KubeconfigPath string
	ContextName    string
}

// Start creates and starts a readonly reverse proxy on a random localhost port.
// The proxy loads TLS/auth config from the kubeconfig and forwards only
// GET, HEAD, and OPTIONS requests (without protocol upgrades) to the real API server.
func Start(cfg Config) (*ReadonlyProxy, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: cfg.KubeconfigPath}
	overrides := &clientcmd.ConfigOverrides{CurrentContext: cfg.ContextName}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)

	restCfg, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	targetURL, err := url.Parse(restCfg.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL %q: %w", restCfg.Host, err)
	}

	transport, err := rest.TransportFor(restCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	handler := NewHandler(targetURL, transport)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	srv := &http.Server{Handler: handler}
	go srv.Serve(listener)

	debugLog.Printf("started on %s, proxying to %s", listener.Addr(), targetURL)

	return &ReadonlyProxy{
		server:   srv,
		listener: listener,
	}, nil
}

// Addr returns the listener address (e.g. "127.0.0.1:54321").
func (p *ReadonlyProxy) Addr() string {
	return p.listener.Addr().String()
}

// Shutdown gracefully stops the proxy.
func (p *ReadonlyProxy) Shutdown(ctx context.Context) error {
	debugLog.Printf("shutting down")
	return p.server.Shutdown(ctx)
}

// NewHandler creates the readonly proxy HTTP handler.
// Exported for testing with a fake backend.
func NewHandler(target *url.URL, transport http.RoundTripper) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = transport
	proxy.FlushInterval = -1 // flush immediately for streaming (logs -f, watches)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		debugLog.Printf(">> %s %s", r.Method, r.URL.Path)

		if reason, ok := checkRequest(r); !ok {
			debugLog.Printf("<< %s %s -> 405 (%s)", r.Method, r.URL.Path, reason)
			writeBlockedResponse(w, r.Method,
				fmt.Sprintf("[kubectx] readonly mode: %s", reason))
			return
		}

		debugLog.Printf("<< %s %s -> proxied", r.Method, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})
}

// checkRequest determines whether a request should be proxied or blocked.
// Returns ("", true) if allowed, or (reason, false) if blocked.
func checkRequest(r *http.Request) (reason string, allowed bool) {
	if isUpgrade(r) {
		return "operations like exec, cp, and port-forward are not allowed", false
	}
	if isReadOnly(r) {
		return "", true
	}
	if isNonMutatingPost(r) {
		return "", true
	}
	if isDryRun(r) {
		return "", true
	}
	return fmt.Sprintf("%s requests are not allowed", r.Method), false
}

// isUpgrade returns true if the request is a protocol upgrade (SPDY/WebSocket).
func isUpgrade(r *http.Request) bool {
	return r.Header.Get("Connection") == "Upgrade" || r.Header.Get("Upgrade") != ""
}

// isReadOnly returns true for safe HTTP methods that never modify state.
func isReadOnly(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}

// isNonMutatingPost returns true for Kubernetes "review" endpoints that use
// POST but don't create persistent resources (e.g. SubjectAccessReview).
// Patterns are anchored to known API groups to prevent spoofing.
func isNonMutatingPost(r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}
	for _, re := range nonMutatingPostPatterns {
		if re.MatchString(r.URL.Path) {
			return true
		}
	}
	return false
}

// isDryRun returns true if the request has ?dryRun=All, which means
// the API server will validate but not persist the request.
func isDryRun(r *http.Request) bool {
	return r.URL.Query().Get("dryRun") == "All"
}

func writeBlockedResponse(w http.ResponseWriter, method, message string) {
	status := &metav1.Status{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Status"},
		Status:   metav1.StatusFailure,
		Message:  message,
		Reason:   metav1.StatusReasonMethodNotAllowed,
		Code:     http.StatusMethodNotAllowed,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(status)
}
