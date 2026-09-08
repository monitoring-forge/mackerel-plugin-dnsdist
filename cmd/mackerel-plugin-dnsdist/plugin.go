package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (opt *Opt) httpClient() *http.Client {
	transport := &http.Transport{
		// inherited http.DefaultTransport
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   opt.Timeout,
			KeepAlive: opt.Timeout,
		}).DialContext,
		TLSHandshakeTimeout:   opt.Timeout,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: opt.Timeout,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   opt.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (opt *Opt) MetricKeyPrefix() string {
	if opt.Prefix == "" {
		opt.Prefix = "dnsdist"
	}
	return opt.Prefix
}

func (opt *Opt) MetricsDefinition(label string, metrics []mp.Metrics) mp.Graphs {
	labelPrefix := cases.Title(language.Und, cases.NoLower).String(opt.Prefix)
	return mp.Graphs{
		Label:   labelPrefix + ": " + label,
		Unit:    "integer",
		Metrics: metrics,
	}
}

func (opt *Opt) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"acl-drop": opt.MetricsDefinition("Dropped packets because of the ACL", []mp.Metrics{
			{Name: "acl-drops", Label: "Dropped", Diff: true},
		}),
		"cache": opt.MetricsDefinition("Packet Cache", []mp.Metrics{
			{Name: "cache-hits", Label: "Hits", Stacked: true, Diff: true},
			{Name: "cache-misses", Label: "Misses", Stacked: true, Diff: true},
		}),
		"downstream-errors": opt.MetricsDefinition("Backend errors", []mp.Metrics{
			{Name: "downstream-send-errors", Label: "Send error", Diff: true},
			{Name: "downstream-timeouts", Label: "Timeouts", Diff: true},
		}),
		"latency": opt.MetricsDefinition("Latency (microseconds)", []mp.Metrics{
			{Name: "latency-avg100", Label: "Latency100"},
			{Name: "latency-avg1000", Label: "Latency1000"},
			{Name: "latency-avg10000", Label: "Latency10000"},
			{Name: "latency-avg1000000", Label: "Latency1000000"},
		}),
		"queries": opt.MetricsDefinition("Queries", []mp.Metrics{
			{Name: "queries", Label: "Queries", Diff: true},
			{Name: "rdqueries", Label: "Query with rd bit", Diff: true},
		}),
		"responses": opt.MetricsDefinition("Response", []mp.Metrics{
			{Name: "responses", Label: "Backend responses", Diff: true},
			{Name: "self-answered", Label: "Self answered", Diff: true},
			{Name: "servfail-responses", Label: "Backend servfail", Diff: true},
		}),
		"rule": opt.MetricsDefinition("Returned because of rules", []mp.Metrics{
			{Name: "rule-drop", Label: "Drop", Stacked: true, Diff: true},
			{Name: "rule-nxdomain", Label: "Nxdomain", Stacked: true, Diff: true},
			{Name: "rule-refused", Label: "Refused", Stacked: true, Diff: true},
			{Name: "rule-servfail", Label: "Servfail", Stacked: true, Diff: true},
			{Name: "rule-truncated", Label: "Truncated", Stacked: true, Diff: true},
		}),
		"fd": opt.MetricsDefinition("FD usage", []mp.Metrics{
			{Name: "fd-usage", Label: "usage"},
		}),
	}
}

func (opt *Opt) FetchMetrics() (map[string]float64, error) {
	req, err := http.NewRequest("GET", opt.getURL(), nil)
	if err != nil {
		return nil, err
	}
	if apiKey := opt.getAPIKey(); apiKey != "" {
		req.Header.Add("X-API-Key", apiKey)
	}
	res, err := opt.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	t := map[string]any{}
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()

	if err := decoder.Decode(&t); err != nil {
		return nil, err
	}

	result := map[string]float64{}
	for k, b := range t {
		if num, ok := b.(json.Number); ok {
			if f, err := num.Float64(); err == nil {
				result[k] = f
			}
		}
	}
	return result, nil
}
