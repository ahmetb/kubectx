// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

const (
	maxHistoryLength = 10
)

// ContextHistory stores the history of used contexts
type ContextHistory struct {
	Contexts []string `json:"contexts"`
}

// Add adds a context to history, moving it to the front if it already exists
func (h *ContextHistory) Add(ctx string) {
	// Remove if exists
	for i, c := range h.Contexts {
		if c == ctx {
			h.Contexts = append(h.Contexts[:i], h.Contexts[i+1:]...)
			break
		}
	}
	
	// Add to front
	h.Contexts = append([]string{ctx}, h.Contexts...)
	
	// Trim if needed
	if len(h.Contexts) > maxHistoryLength {
		h.Contexts = h.Contexts[:maxHistoryLength]
	}
}

// getHistoryFilePath returns the path to the history file
var getHistoryFilePath = func() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kube", "kubectx_history.json"), nil
}

// loadHistory loads the context history from disk
func loadHistory() (*ContextHistory, error) {
	path, err := getHistoryFilePath()
	if err != nil {
		return nil, err
	}
	
	// If file doesn't exist, return empty history
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &ContextHistory{Contexts: []string{}}, nil
	}
	
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	// Parse JSON
	var history ContextHistory
	if len(data) > 0 {
		if err := json.Unmarshal(data, &history); err != nil {
			return nil, err
		}
	} else {
		history.Contexts = []string{}
	}
	
	return &history, nil
}

// saveHistory saves the context history to disk
func saveHistory(history *ContextHistory) error {
	path, err := getHistoryFilePath()
	if err != nil {
		return err
	}
	
	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Serialize to JSON
	data, err := json.Marshal(history)
	if err != nil {
		return err
	}
	
	// Write to file
	return os.WriteFile(path, data, 0644)
}

// prioritizeContexts sorts contexts with recently used ones first
func prioritizeContexts(allContexts []string, historyContexts []string) []string {
	// Create a map for O(1) lookup
	seen := make(map[string]bool)
	result := []string{}
	
	// First add contexts from history that exist in allContexts
	for _, ctx := range historyContexts {
		if contains(allContexts, ctx) {
			result = append(result, ctx)
			seen[ctx] = true
		}
	}
	
	// Then add remaining contexts in alphabetical order
	remaining := []string{}
	for _, ctx := range allContexts {
		if !seen[ctx] {
			remaining = append(remaining, ctx)
		}
	}
	sort.Strings(remaining)
	
	return append(result, remaining...)
}

// contains checks if a string is in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
