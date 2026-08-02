package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Plugin struct {
	Prefix  string
	URL     string
	Timeout time.Duration
	APIKey  string
}

func (p *Plugin) httpClient() *http.Client {
	transport := &http.Transport{
		// inherited http.DefaultTransport
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   p.Timeout,
			KeepAlive: p.Timeout,
		}).DialContext,
		TLSHandshakeTimeout:   p.Timeout,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: p.Timeout,
	}
	return &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (p *Plugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "dnsdist"
	}
	return p.Prefix
}

func (p *Plugin) MetricsDefinition(label string, metrics []mp.Metrics) mp.Graphs {
	labelPrefix := cases.Title(language.Und, cases.NoLower).String(p.Prefix)
	return mp.Graphs{
		Label:   labelPrefix + ": " + label,
		Unit:    "integer",
		Metrics: metrics,
	}
}

func (p *Plugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"acl-drop": p.MetricsDefinition("Dropped packets becaused of the ACL", []mp.Metrics{
			{Name: "acl-drops", Label: "Dropped", Diff: true},
		}),
		"cache": p.MetricsDefinition("Packet Cache", []mp.Metrics{
			{Name: "cache-hits", Label: "Hits", Stacked: true, Diff: true},
			{Name: "cache-misses", Label: "Misses", Stacked: true, Diff: true},
		}),
		"downstream-errors": p.MetricsDefinition("Backend errors", []mp.Metrics{
			{Name: "downstream-send-errors", Label: "Send error", Diff: true},
			{Name: "downstream-timeouts", Label: "Timeouts", Diff: true},
		}),
		"latency": p.MetricsDefinition("Latency (microseconds)", []mp.Metrics{
			{Name: "latency-avg100", Label: "Latency100"},
			{Name: "latency-avg1000", Label: "Latency1000"},
			{Name: "latency-avg10000", Label: "Latency10000"},
			{Name: "latency-avg1000000", Label: "Latency1000000"},
		}),
		"queries": p.MetricsDefinition("Queries", []mp.Metrics{
			{Name: "queries", Label: "Queries", Diff: true},
			{Name: "rdqueries", Label: "Query with rd bit", Diff: true},
		}),
		"responses": p.MetricsDefinition("Response", []mp.Metrics{
			{Name: "responses", Label: "Backend responses", Diff: true},
			{Name: "self-answered", Label: "Self answered", Diff: true},
			{Name: "servfail-responses", Label: "Backend servfail", Diff: true},
		}),
		"rule": p.MetricsDefinition("Returned because of rules", []mp.Metrics{
			{Name: "rule-drop", Label: "Drop", Stacked: true, Diff: true},
			{Name: "rule-nxdomain", Label: "Nxdomain", Stacked: true, Diff: true},
			{Name: "rule-refused", Label: "Refused", Stacked: true, Diff: true},
			{Name: "rule-servfail", Label: "Servfail", Stacked: true, Diff: true},
			{Name: "rule-truncated", Label: "Truncated", Stacked: true, Diff: true},
		}),
		"fd": p.MetricsDefinition("FD usage", []mp.Metrics{
			{Name: "fd-usage", Label: "usage"},
		}),
	}
}

func (p *Plugin) FetchMetrics() (map[string]float64, error) {
	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		return nil, err
	}
	if p.APIKey != "" {
		req.Header.Add("X-API-Key", p.APIKey)
	}
	res, err := p.httpClient().Do(req)
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
		f, err := strconv.ParseFloat(fmt.Sprintf("%v", b), 64)
		if err != nil {
			continue
		}
		result[k] = f
	}
	return result, nil
}

func (u *Plugin) Run() {
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
}
