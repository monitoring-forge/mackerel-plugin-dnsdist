# mackerel-plugin-dnsdist

Mackerel プラグイン for dnsdist。dnsdist は DNS、DoS、不正アクセスを高度に意識したロードバランサーです。

https://dnsdist.org/

## インストール

リリースページからダウンロードするか、Mackerel プラグインレジストリからインストールしてください:

```sh
mkr plugin install monitoring-forge/mackerel-plugin-dnsdist
```

## 動作要件

- 内蔵 Web サーバーが有効化され、API キーが設定された dnsdist。
- このプラグインは dnsdist Web サーバーに対して `GET /jsonstat?command=stats` を問い合わせます。
  詳細は [dnsdist Web サーバードキュメント](https://www.dnsdist.org/guides/webserver.html#get--jsonstat-query-parameters) を参照してください。

## 使い方

```sh
mackerel-plugin-dnsdist [OPTIONS]
```

### オプション

| オプション | デフォルト | 説明 |
| --- | --- | --- |
| `-v`, `--version` | - | バージョン情報を表示する |
| `-p`, `--port` | `8083` | dnsdist Web サーバーのポート番号 |
| `-H`, `--hostname` | `127.0.0.1` | dnsdist Web サーバーのホスト名または IP アドレス |
| `--prefix` | `dnsdist` | Mackerel で使用するメトリックキーのプレフィックス |
| `--timeout` | `30s` | HTTP リクエストのタイムアウト |
| `--api-key` | - | dnsdist Web サーバーの API キー（X-API-Key ヘッダー） |

## 認証

このプラグインは API キーを使用して dnsdist Web サーバーに対して認証を行います。
キーは `X-API-Key` リクエストヘッダーに含めて送信されます。

API キーは以下の順序で解決されます:

1. `--api-key` オプションで指定された値。
2. 環境変数 `DNSDIST_CONFIG_PATH` で指定されたファイルから読み取られた値。
3. `/etc/dnsdist/dnsdist.conf` から読み取られた値。

dnsdist 設定ファイルから読み取る場合、プラグインは以下のような行からキーを抽出します:

```lua
setWebserverConfig(..., { ..., apiKey = "supersecretAPIkey", ... })
```

## メトリクス

このプラグインは `/jsonstat?command=stats` から以下のメトリクスを収集します。
dnsdist が返す統計情報の全一覧については、[dnsdist 統計情報ドキュメント](https://www.dnsdist.org/statistics.html) を参照してください。

| グラフ | メトリック | 説明 |
| --- | --- | --- |
| acl-drop | acl-drops | ACL によってドロップされたパケット数 |
| cache | cache-hits | キャッシュから回答が取得された回数 |
| cache | cache-misses | キャッシュに回答が見つからなかった回数 |
| downstream-errors | downstream-send-errors | バックエンドにクエリを送信した際のエラー数 |
| downstream-errors | downstream-timeouts | バックエンドから時間内に回答が返ってこなかったクエリ数 |
| latency | latency-avg100 | 直近 100 パケットの平均応答レイテンシー（マイクロ秒） |
| latency | latency-avg1000 | 直近 1000 パケットの平均応答レイテンシー（マイクロ秒） |
| latency | latency-avg10000 | 直近 10000 パケットの平均応答レイテンシー（マイクロ秒） |
| latency | latency-avg1000000 | 直近 1000000 パケットの平均応答レイテンシー（マイクロ秒） |
| queries | queries | 受信したクエリ数 |
| queries | rdqueries | 再帰希望ビットが設定された受信クエリ数 |
| responses | responses | バックエンドから受信した応答数 |
| responses | self-answered | 自己応答した応答数 |
| responses | servfail-responses | バックエンドから受信した SERVFAIL 回答数 |
| rule | rule-drop | ルールによってドロップされたクエリ数 |
| rule | rule-nxdomain | ルールによって返された NXDomain 回答数 |
| rule | rule-refused | ルールによって返された Refused 回答数 |
| rule | rule-servfail | ルールによって返された SERVFAIL 回答数 |
| rule | rule-truncated | ルールによって返された切り詰め応答数 |
| fd | fd-usage | 現在使用中のファイルディスクリプタ数 |

## 例

```sh
mackerel-plugin-dnsdist -H 127.0.0.1 -p 8083 --api-key supersecretAPIkey
```

カスタムメトリックプレフィックスを指定する場合:

```sh
mackerel-plugin-dnsdist -H 127.0.0.1 -p 8083 --prefix dnsdist-prod
```

## ライセンス

[LICENSE](LICENSE) を参照してください。
