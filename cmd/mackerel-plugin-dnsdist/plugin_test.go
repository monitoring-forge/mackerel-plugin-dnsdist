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
		// Response sample from https://www.dnsdist.org/guides/webserver.html#get--jsonstat
		_, err := fmt.Fprintln(w, `{"acl-drops": 0, "cache-hits": 0, "cache-misses": 0, "cpu-sys-msec": 633, "cpu-user-msec": 499, "downstream-send-errors": 0, "downstream-timeouts": 0, "dyn-block-nmg-size": 1, "dyn-blocked": 3, "empty-queries": 0, "fd-usage": 17, "latency-avg100": 7651.3982737482893, "latency-avg1000": 860.05142763680249, "latency-avg10000": 87.032142373878372, "latency-avg1000000": 0.87146026426551759, "latency-slow": 0, "latency0-1": 0, "latency1-10": 0, "latency10-50": 22, "latency100-1000": 1, "latency50-100": 0, "no-policy": 0, "noncompliant-queries": 0, "noncompliant-responses": 0, "over-capacity-drops": 0, "packetcache-hits": 0, "packetcache-misses": 0, "queries": 26, "rdqueries": 26, "real-memory-usage": 6078464, "responses": 23, "rule-drop": 0, "rule-nxdomain": 0, "rule-refused": 0, "self-answered": 0, "server-policy": "leastOutstanding", "servfail-responses": 0, "too-old-drops": 0, "trunc-failures": 0, "uptime": 412, "dummy-text": "dummy-text-value"}`)
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
	assert.Equal(t, 7651.3982737482893, metrics["latency-avg100"], "Expected latency-avg100 metric to be 7651.3982737482893")
	assert.Equal(t, 17.0, metrics["fd-usage"], "Expected fd-usage metric to be 17")
	assert.NotContains(t, metrics, "dummy-text", "Expected metrics not to contain dummy-text")
}
