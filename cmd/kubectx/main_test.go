package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextHistory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "kubectx-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override history file path for testing
	origHistoryFilePath := getHistoryFilePath
	defer func() { getHistoryFilePath = origHistoryFilePath }()
	
	historyPath := filepath.Join(tmpDir, "history.json")
	getHistoryFilePath = func() (string, error) {
		return historyPath, nil
	}

	t.Run("New history is empty", func(t *testing.T) {
		history, err := loadHistory()
		assert.NoError(t, err)
		assert.NotNil(t, history)
		assert.Empty(t, history.Contexts)
	})

	t.Run("Add context to history", func(t *testing.T) {
		history := &ContextHistory{Contexts: []string{}}
		history.Add("context1")
		assert.Equal(t, []string{"context1"}, history.Contexts)
		
		// Add the same context again
		history.Add("context1")
		assert.Equal(t, []string{"context1"}, history.Contexts)
		
		// Add another context
		history.Add("context2")
		assert.Equal(t, []string{"context2", "context1"}, history.Contexts)
	})

	t.Run("Save and load history", func(t *testing.T) {
		history := &ContextHistory{Contexts: []string{"context1", "context2"}}
		err := saveHistory(history)
		assert.NoError(t, err)
		
		loaded, err := loadHistory()
		assert.NoError(t, err)
		assert.Equal(t, history.Contexts, loaded.Contexts)
	})

	t.Run("History respects max length", func(t *testing.T) {
		history := &ContextHistory{Contexts: []string{}}
		
		// Add more than maxHistoryLength contexts
		for i := 0; i < maxHistoryLength+5; i++ {
			history.Add(fmt.Sprintf("context%d", i))
		}
		
		assert.Equal(t, maxHistoryLength, len(history.Contexts))
		assert.Equal(t, "context9", history.Contexts[0])
	})
}

func TestPrioritizeContexts(t *testing.T) {
	allContexts := []string{"a", "b", "c", "d", "e"}
	historyContexts := []string{"c", "a", "f"} // Note: "f" is not in allContexts
	
	result := prioritizeContexts(allContexts, historyContexts)
	
	// Expected: c, a (from history, in order), then b, d, e alphabetically
	expected := []string{"c", "a", "b", "d", "e"}
	assert.Equal(t, expected, result)
}
