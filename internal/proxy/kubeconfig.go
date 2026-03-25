package proxy

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// RewriteKubeconfig takes minified kubeconfig bytes and rewrites them so that:
//   - The cluster server URL points to the local proxy address (plain HTTP).
//   - insecure-skip-tls-verify is set (needed for plain HTTP).
//   - Certificate authority data is removed.
//   - User auth fields (client certs, tokens, exec, auth-provider) are removed
//     since the proxy handles authentication to the real API server.
func RewriteKubeconfig(data []byte, proxyAddr string) ([]byte, error) {
	cfg, err := clientcmd.Load(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	for _, cluster := range cfg.Clusters {
		cluster.Server = "http://" + proxyAddr
		cluster.InsecureSkipTLSVerify = true
		cluster.CertificateAuthority = ""
		cluster.CertificateAuthorityData = nil
	}

	for name := range cfg.AuthInfos {
		cfg.AuthInfos[name] = &clientcmdapi.AuthInfo{}
	}

	// Rename contexts with [RO] suffix to indicate readonly mode.
	renames := make(map[string]string, len(cfg.Contexts))
	for name := range cfg.Contexts {
		renames[name] = name + "[RO]"
	}
	for old, roName := range renames {
		cfg.Contexts[roName] = cfg.Contexts[old]
		delete(cfg.Contexts, old)
		if cfg.CurrentContext == old {
			cfg.CurrentContext = roName
		}
	}

	out, err := clientcmd.Write(*cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize kubeconfig: %w", err)
	}
	return out, nil
}
