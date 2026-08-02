package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchMetrics(t *testing.T) {
	expectedAPIKey := "test-api-key"

	// Create a test server to simulate the dnsdist API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != expectedAPIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprintln(w, `{"rule-drop": 10, "rule-nxdomain": 5}`)
		if err != nil {
			t.Logf("Failed to write response: %v", err)
		}
	}))
	defer ts.Close()

	plugin := &Plugin{
		URL:    ts.URL,
		APIKey: expectedAPIKey,
	}

	metrics, err := plugin.FetchMetrics()
	assert.NoError(t, err, "Expected no error from FetchMetrics")
	assert.Equal(t, 10.0, metrics["rule-drop"], "Expected rule-drop metric to be 10")
	assert.Equal(t, 5.0, metrics["rule-nxdomain"], "Expected rule-nxdomain metric to be 5")
}
