package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/jessevdk/go-flags"
)

var version string
var commit string

const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

type Opt struct {
	Version bool   `short:"v" long:"version" description:"Show version"`
	Prefix  string `long:"prefix" default:"dnsdist" description:"Metric key prefix"`

	Port    string        `short:"p" long:"port" default:"8083" description:"Port number"`
	Host    string        `short:"H" long:"hostname" default:"127.0.0.1" description:"Hostname"`
	Timeout time.Duration `long:"timeout" default:"30s" description:"Timeout"`

	APIKey string `long:"api-key" description:"api key"`
}

func (o *Opt) URL() string {
	url := url.URL{
		Scheme:   "http",
		Host:     net.JoinHostPort(o.Host, o.Port),
		Path:     "/jsonstat",
		RawQuery: "command=stats",
	}
	return url.String()
}

var apiKeyRegexp = regexp.MustCompile(`setWebserverConfig\(.*\{.*\bapiKey\s*=\s*"(.+?)"`)

func (o *Opt) GetAPIKey() string {
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

func main() {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		os.Exit(OK)
	} else if flags.WroteHelp(err) {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(OK)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(UNKNOWN)
	}

	u := &Plugin{
		Prefix:  opt.Prefix,
		Timeout: opt.Timeout,
		URL:     opt.URL(),
		APIKey:  opt.GetAPIKey(),
	}
	u.Run()
}
