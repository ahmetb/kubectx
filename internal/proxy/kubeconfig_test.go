package proxy

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

const sampleKubeconfig = `apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: dGVzdC1jYS1kYXRh
    server: https://my-cluster.example.com:6443
  name: my-cluster
contexts:
- context:
    cluster: my-cluster
    user: my-user
  name: my-context
current-context: my-context
users:
- name: my-user
  user:
    client-certificate-data: dGVzdC1jZXJ0LWRhdGE=
    client-key-data: dGVzdC1rZXktZGF0YQ==
    token: test-token
`

func TestRewriteKubeconfig(t *testing.T) {
	result, err := RewriteKubeconfig([]byte(sampleKubeconfig), "127.0.0.1:12345")
	if err != nil {
		t.Fatalf("RewriteKubeconfig failed: %v", err)
	}

	cfg, err := clientcmd.Load(result)
	if err != nil {
		t.Fatalf("failed to parse rewritten kubeconfig: %v", err)
	}

	cluster := cfg.Clusters["my-cluster"]
	if cluster == nil {
		t.Fatal("cluster my-cluster not found")
	}
	if cluster.Server != "http://127.0.0.1:12345" {
		t.Errorf("expected server http://127.0.0.1:12345, got %q", cluster.Server)
	}
	if !cluster.InsecureSkipTLSVerify {
		t.Error("expected insecure-skip-tls-verify to be true")
	}
	if len(cluster.CertificateAuthorityData) != 0 {
		t.Error("expected certificate-authority-data to be removed")
	}
	if cluster.CertificateAuthority != "" {
		t.Error("expected certificate-authority to be removed")
	}

	user := cfg.AuthInfos["my-user"]
	if user == nil {
		t.Fatal("user my-user not found")
	}
	if len(user.ClientCertificateData) != 0 {
		t.Error("expected client-certificate-data to be removed")
	}
	if len(user.ClientKeyData) != 0 {
		t.Error("expected client-key-data to be removed")
	}
	if user.Token != "" {
		t.Error("expected token to be removed")
	}

	if _, ok := cfg.Contexts["my-context[RO]"]; !ok {
		t.Error("expected context to be renamed to \"my-context[RO]\"")
	}
	if _, ok := cfg.Contexts["my-context"]; ok {
		t.Error("expected original context name to be removed")
	}
	if cfg.CurrentContext != "my-context[RO]" {
		t.Errorf("expected current-context to be \"my-context[RO]\", got %q", cfg.CurrentContext)
	}
}

const execKubeconfig = `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://my-cluster.example.com:6443
  name: my-cluster
contexts:
- context:
    cluster: my-cluster
    user: my-user
  name: my-context
current-context: my-context
users:
- name: my-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: gke-gcloud-auth-plugin
`

func TestRewriteKubeconfig_ExecPlugin(t *testing.T) {
	result, err := RewriteKubeconfig([]byte(execKubeconfig), "127.0.0.1:54321")
	if err != nil {
		t.Fatalf("RewriteKubeconfig failed: %v", err)
	}

	cfg, err := clientcmd.Load(result)
	if err != nil {
		t.Fatalf("failed to parse rewritten kubeconfig: %v", err)
	}

	user := cfg.AuthInfos["my-user"]
	if user == nil {
		t.Fatal("user my-user not found")
	}
	if user.Exec != nil {
		t.Error("expected exec plugin config to be removed")
	}
}
