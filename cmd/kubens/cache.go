package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type kubeNsCache map[string][]string

var cacheRelativePath = []string{".kube", "kubens-cache"}

func newKubensCache() (kubeNsCache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to init kubens cache, %w", err)
	}

	cachePath := filepath.Join(home, filepath.Join(cacheRelativePath...))

	cacheHandle, err := os.Open(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache file: %w", err)
	}

	defer cacheHandle.Close()

	cache := make(kubeNsCache)

	if err := json.NewDecoder(cacheHandle).Decode(&cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	return cache, nil
}

func (c *kubeNsCache) Save(ctx string, namespaces []string) error {
	(*c)[ctx] = namespaces

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory")
	}

	cachePath := filepath.Join(home, filepath.Join(cacheRelativePath...))

	cacheHandle, err := os.Create(cachePath)
	if err != nil {
		return fmt.Errorf("failed to create cache file")
	}

	defer cacheHandle.Close()

	err = json.NewEncoder(cacheHandle).Encode(c)
	if err != nil {
		return fmt.Errorf("failed to encode cache file")
	}

	return nil
}

func (c *kubeNsCache) Get(ctx string) ([]string, bool) {
	namespaces, ok := (*c)[ctx]
	if !ok {
		return nil, false
	}
	return namespaces, true
}
