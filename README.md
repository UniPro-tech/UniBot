# UniBot

[![CodeQL](https://github.com/UniPro-tech/UniBot/actions/workflows/github-code-scanning/codeql/badge.svg?branch=main)](https://github.com/UniPro-tech/UniBot/actions/workflows/github-code-scanning/codeql) [![Docker Build and Push](https://github.com/UniPro-tech/UniBot/actions/workflows/docker-tag.yaml/badge.svg)](https://github.com/UniPro-tech/UniBot/actions/workflows/docker-tag.yaml)

UniProjectのDiscordサーバーで運用されているDiscordBot

## Environment Variables

環境変数は下記の通りです。

- `DISCORD_TOKEN` - Discord Token
- `CONFIG_ADMIN_GUILD_ID` - Admin Guild ID
- `CONFIG_ADMIN_ROLE_ID` - Admin Role ID
- `PG_DSN` - DSN for PostgreSQL connection
- `VOICEVOX_URI` - (Optional) VOICEVOX Engine URI. Default: localhost:53000
- `VOICEVOX_API_KEY` - (Optional) API Key passed as Authorization: ApiKey when communicating with VOICEVOX Engine

### Logging

ログはstdoutにJSONとして出力されます。Kubernetes上のfluentbitにより、Grafana / Lokiに収集されることを想定しています。
全てのレコードには `trace_id` と `request_id` が付帯しており、単一のDiscordイベントをgoroutineにまたがって追跡できるようになっています。

- `CONFIG_LOG_LEVEL` - (Optional) `debug` / `info` / `notice` / `warn` / `error`. Default: `info`
- `CONFIG_LOG_FORMAT` - (Optional) `json` or `text`. Use `text` for local development. Default: `json`
- `CONFIG_LOG_SOURCE` - (Optional) `true` to include the source file and line of each record. Default: `false`
- `CONFIG_LOG_SQL` - (Optional) `silent` / `error` / `warn` / `info` for gorm. `info` logs every query at debug level. Default: `warn`
- `CONFIG_LOG_SQL_SLOW_MS` - (Optional) Queries slower than this are logged as `slow sql`. Default: `200`
- `CONFIG_LOG_ERROR_CHANNEL_ID` - (Optional) Discord channel to notify about Error / Warn / Notice logs. **Unset disables Discord notification entirely.**
- `CONFIG_LOG_DISCORD_LEVEL` - (Optional) Minimum level forwarded to Discord. Default: `notice`
- `CONFIG_LOG_READY_CHANNEL_ID` - (Optional) Discord channel for startup / shutdown notices. Falls back to `CONFIG_LOG_ERROR_CHANNEL_ID`

Discordの通知には意図的に `trace_id`、 `request_id`、 `level`、 `timestamp` **のみ** が含まれており、ログメッセージ、エラー文、その他の情報は一切含みません。
ログ本文はユーザーのメッセージ内容、辞書の中身、そしてSQLが含まれている可能性があるため、代わりに詳細はGrafanaで `trace_id` を用いて見つけることができます。

> [!WARNING]
> `CONFIG_LOG_LEVEL=debug` はdisgoにBotが受け取った、もしくは送信したメッセージの内容が含まれた **全てのRESTリクエストおよびレスポンス** のログを出力させます。
> 本番環境では使用しないでください。

## Running with Docker

ローカルでビルドすることで、Dockerを使って起動することもできます。

### Using docker-compose

`_docker-compose.prod.yaml` を `docker-compose.yaml` に変更し、必要事項を記述してください。

## Building and Developing Yourself

### Prerequisites

下記の依存関係のインストールが必要です。
下記はDockerfileから確認できる依存関係です。実環境にインストールする際は、OSに応じてそれぞれのインストール方法を確認してください。

- Go >= 1.24
- opus
- opus-dev
- opusfile
- opusfile-dev
- ffmpeg

### Configuration

`make run` を行うことで起動できます。

### About Database

The database uses PostgreSQL, and you can start only the database using `docker-compose.dev.yaml`.

### Installing Go Dependencies

```bash
go mod tidy
```

### Running Only

```bash
make run
```

or

```bash
make run-rss
```

### Building

```bash
make build
```

or

```bash
make build-rss
```

## Built With

- [disgo](https://pkg.go.dev/github.com/disgoorg/disgo) - The Discord API wrapper for Golang.
- [ohraban/opus](https://pkg.go.dev/github.com/hraban/opus) - The Golang bindings for the xiph.org C libraries libopus and libopusfile.
- [gorm](https://gorm.io/ja_JP/) - The fantastic ORM library for Golang.

## Contributing

[CONTRIBUTING.md](CONTRIBUTING.md) をお読みください。
Code of Conductの詳細と、Pull Requestの送り方などが記載されています。

## Versioning

[Semantic Versioning 2.0](https://semver.org/lang/ja/) を使用しております。
また、それぞれのコンポーネントに対して下記のようなルールを設けて運用しております。

- bot/vx.x.x
- rss/vx.x.x

## License

このプロジェクトはMITライセンスで保護されています。
詳しくは、[LICENSE.md](LICENSE.md) をご覧ください。
