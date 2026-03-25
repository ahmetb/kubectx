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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ahmetb/kubectx/internal/env"
)

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

		// Block protocol upgrades (exec, cp, port-forward use SPDY/WebSocket).
		if r.Header.Get("Connection") == "Upgrade" || r.Header.Get("Upgrade") != "" {
			debugLog.Printf("<< %s %s -> 405 (upgrade blocked)", r.Method, r.URL.Path)
			writeBlockedResponse(w, r.Method, "[kubectx] readonly mode: operations like exec, cp, and port-forward are not allowed")
			return
		}

		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			debugLog.Printf("<< %s %s -> proxied", r.Method, r.URL.Path)
			proxy.ServeHTTP(w, r)
		default:
			debugLog.Printf("<< %s %s -> 405 (blocked)", r.Method, r.URL.Path)
			writeBlockedResponse(w, r.Method,
				fmt.Sprintf("[kubectx] readonly mode: %s requests are not allowed", r.Method))
		}
	})
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
