# mackerel-plugin-dnsdist

Mackerel Plugin for dnsdist, a highly DNS-, DoS- and abuse-aware load balancer.

https://dnsdist.org/

## Install

Download from the release page or install via the Mackerel plugin registry:

```sh
mkr plugin install monitoring-forge/mackerel-plugin-dnsdist
```

## Requirements

- dnsdist with the built-in webserver enabled and API key configured.
- The plugin queries `GET /jsonstat?command=stats` on the dnsdist webserver.
  See the [dnsdist webserver documentation](https://www.dnsdist.org/guides/webserver.html#get--jsonstat-query-parameters) for details.

## Usage

```sh
mackerel-plugin-dnsdist [OPTIONS]
```

### Options

| Option | Default | Description |
| --- | --- | --- |
| `-v`, `--version` | - | Show version information |
| `-p`, `--port` | `8083` | Port number of the dnsdist webserver |
| `-H`, `--hostname` | `127.0.0.1` | Hostname or IP address of the dnsdist webserver |
| `--prefix` | `dnsdist` | Metric key prefix used in Mackerel |
| `--timeout` | `30s` | HTTP request timeout |
| `--api-key` | - | API key for the dnsdist webserver (X-API-Key header) |

## Authentication

The plugin uses the API key to authenticate against the dnsdist webserver.
The key is sent in the `X-API-Key` request header.

The API key is resolved in the following order:

1. The value provided via the `--api-key` option.
2. The value read from the file specified by the `DNSDIST_CONFIG_PATH` environment variable.
3. The value read from `/etc/dnsdist/dnsdist.conf`.

When reading from a dnsdist configuration file, the plugin extracts the key from a line like:

```lua
setWebserverConfig(..., { ..., apiKey = "supersecretAPIkey", ... })
```

## Metrics

This plugin collects the following metrics from `/jsonstat?command=stats`.
For the full list of statistics returned by dnsdist, refer to the [dnsdist statistics documentation](https://www.dnsdist.org/statistics.html).

| Graph | Metric | Description |
| --- | --- | --- |
| acl-drop | acl-drops | Number of packets dropped because of the ACL |
| cache | cache-hits | Number of times an answer was retrieved from cache |
| cache | cache-misses | Number of times an answer was not found in cache |
| downstream-errors | downstream-send-errors | Number of errors when sending a query to a backend |
| downstream-errors | downstream-timeouts | Number of queries not answered in time by a backend |
| latency | latency-avg100 | Average response latency in microseconds of the last 100 packets |
| latency | latency-avg1000 | Average response latency in microseconds of the last 1000 packets |
| latency | latency-avg10000 | Average response latency in microseconds of the last 10000 packets |
| latency | latency-avg1000000 | Average response latency in microseconds of the last 1000000 packets |
| queries | queries | Number of received queries |
| queries | rdqueries | Number of received queries with the recursion desired bit set |
| responses | responses | Number of responses received from backends |
| responses | self-answered | Number of self-answered responses |
| responses | servfail-responses | Number of SERVFAIL answers received from backends |
| rule | rule-drop | Number of queries dropped because of a rule |
| rule | rule-nxdomain | Number of NXDomain answers returned because of a rule |
| rule | rule-refused | Number of Refused answers returned because of a rule |
| rule | rule-servfail | Number of SERVFAIL answers returned because of a rule |
| rule | rule-truncated | Number of truncated answers returned because of a rule |
| fd | fd-usage | Number of currently used file descriptors |

## Example

```sh
mackerel-plugin-dnsdist -H 127.0.0.1 -p 8083 --api-key supersecretAPIkey
```

Or with a custom metric prefix:

```sh
mackerel-plugin-dnsdist -H 127.0.0.1 -p 8083 --prefix dnsdist-prod
```

## License

See [LICENSE](LICENSE).
