package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIKey(t *testing.T) {
	// Create a temporary dnsdist.conf file
	tmpdir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpdir, "dnsdist.conf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	// Write a sample configuration with an API key
	configContent := `setWebserverConfig({apiKey = "test-api-key"})`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_ = os.Setenv("DNSDIST_CONFIG_PATH", tmpFile.Name())
	defer func() {
		_ = os.Unsetenv("DNSDIST_CONFIG_PATH")
	}()

	{
		opt := &Opt{}
		apiKey := opt.GetAPIKey()
		expectedAPIKey := "test-api-key"
		assert.Equal(t, expectedAPIKey, apiKey, "Expected API key %s, but got %s", expectedAPIKey, apiKey)
	}

	{
		opt := &Opt{APIKey: "explicit-api-key"}
		apiKey := opt.GetAPIKey()
		expectedAPIKey := "explicit-api-key"
		assert.Equal(t, expectedAPIKey, apiKey, "Expected API key %s, but got %s", expectedAPIKey, apiKey)
	}
}

func TestURL(t *testing.T) {
	opt := &Opt{
		Host: "localhost",
		Port: "8083",
	}
	expectedURL := "http://localhost:8083/jsonstat?command=stats"
	assert.Equal(t, expectedURL, opt.URL(), "Expected URL %s, but got %s", expectedURL, opt.URL())
}
