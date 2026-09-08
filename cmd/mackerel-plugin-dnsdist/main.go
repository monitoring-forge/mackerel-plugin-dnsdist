package main

import (
	"net"
	"net/url"
	"os"
	"regexp"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/monitoring-forge/flagrun"
)

var version string

type Opt struct {
	Version bool   `short:"v" long:"version" description:"Show version"`
	Prefix  string `long:"prefix" default:"dnsdist" description:"Metric key prefix"`

	Port    string        `short:"p" long:"port" default:"8083" description:"Port number"`
	Host    string        `short:"H" long:"hostname" default:"127.0.0.1" description:"Hostname"`
	Timeout time.Duration `long:"timeout" default:"30s" description:"Timeout"`

	APIKey  string `long:"api-key" description:"api key"`
	testURL string // for testing purposes, allows overriding the URL
}

func (o *Opt) getURL() string {
	if o.testURL != "" {
		return o.testURL
	}
	url := url.URL{
		Scheme:   "http",
		Host:     net.JoinHostPort(o.Host, o.Port),
		Path:     "/jsonstat",
		RawQuery: "command=stats",
	}
	return url.String()
}

var apiKeyRegexp = regexp.MustCompile(`setWebserverConfig\(.*\{.*\bapiKey\s*=\s*"(.+?)"`)

func (o *Opt) getAPIKey() string {
	if o.APIKey != "" {
		return o.APIKey
	}

	if configPath := os.Getenv("DNSDIST_CONFIG_PATH"); configPath != "" {
		if apiKey := getAPIKeyFromFile(configPath); apiKey != "" {
			return apiKey
		}
	}

	return getAPIKeyFromFile("/etc/dnsdist/dnsdist.conf")
}

func getAPIKeyFromFile(path string) string {
	buf, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	res := apiKeyRegexp.FindAllSubmatch(buf, -1)
	if len(res) < 1 {
		return ""
	}
	return string(res[0][1])
}

func (opt *Opt) Run(_ []string) {
	plugin := mp.NewMackerelPlugin(opt)
	plugin.Run()
}

func main() {
	opt := &Opt{}
	os.Exit(flagrun.Ship(opt, flagrun.Version(version)))
}
