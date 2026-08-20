# UniBot

[![CodeQL](https://github.com/UniPro-tech/UniBot/actions/workflows/github-code-scanning/codeql/badge.svg?branch=main)](https://github.com/UniPro-tech/UniBot/actions/workflows/github-code-scanning/codeql) [![Docker Build and Push](https://github.com/UniPro-tech/UniBot/actions/workflows/docker-tag.yaml/badge.svg)](https://github.com/UniPro-tech/UniBot/actions/workflows/docker-tag.yaml)

A Discord Bot that manages and operates within All-Japan Digital Creative Club UniProject.

## Environment Variables

- `DISCORD_TOKEN` - Discord Token
- `CONFIG_ADMIN_GUILD_ID` - Admin Guild ID
- `CONFIG_ADMIN_ROLE_ID` - Admin Role ID
- `PG_DSN` - DSN for PostgreSQL connection
- `VOICEVOX_URI` - (Optional) VOICEVOX Engine URI. Default: localhost:53000
- `VOICEVOX_API_KEY` - (Optional) API Key passed as Authorization: ApiKey when communicating with VOICEVOX Engine

### Logging

Logs are written to stdout as one JSON object per line, intended to be collected by
Kubernetes and shipped to Grafana / Loki. Every record carries a `trace_id` and
`request_id` so a single Discord event can be followed across goroutines.

- `CONFIG_LOG_LEVEL` - (Optional) `debug` / `info` / `notice` / `warn` / `error`. Default: `info`
- `CONFIG_LOG_FORMAT` - (Optional) `json` or `text`. Use `text` for local development. Default: `json`
- `CONFIG_LOG_SOURCE` - (Optional) `true` to include the source file and line of each record. Default: `false`
- `CONFIG_LOG_SQL` - (Optional) `silent` / `error` / `warn` / `info` for gorm. `info` logs every query at debug level. Default: `warn`
- `CONFIG_LOG_SQL_SLOW_MS` - (Optional) Queries slower than this are logged as `slow sql`. Default: `200`
- `CONFIG_LOG_ERROR_CHANNEL_ID` - (Optional) Discord channel to notify about Error / Warn / Notice logs. **Unset disables Discord notification entirely.**
- `CONFIG_LOG_DISCORD_LEVEL` - (Optional) Minimum level forwarded to Discord. Default: `notice`
- `CONFIG_LOG_READY_CHANNEL_ID` - (Optional) Discord channel for startup / shutdown notices. Falls back to `CONFIG_LOG_ERROR_CHANNEL_ID`

Discord notifications intentionally contain **only** `trace_id`, `request_id`, `level` and
a timestamp - never the log message, the error text, or any attribute. Log bodies can
contain user message content, dictionary words and SQL, so the details are looked up in
Grafana by `trace_id` instead.

> [!WARNING]
> `CONFIG_LOG_LEVEL=debug` makes disgo log the **body of every REST request and response**,
> which includes the content of messages the bot sends and receives.
> Never use it in production.

## Running with Docker

You can use Docker Image, but you need to build it locally.

### Using docker-compose

Rename `_docker-compose.prod.yaml` to `docker-compose.yaml`, fill in the necessary information, and build it.

## Building and Developing Yourself

### Prerequisites

You need to install these dependencies.
This is an excerpt from the Dockerfile, so please check the installation method for each OS and environment on your own.

- Go >= 1.24
- opus
- opus-dev
- opusfile
- opusfile-dev
- ffmpeg

### Configuration

> [!TIP]
> If you have trouble setting environment variables, try adding export to the relevant parts in the shell script.

Rename `scripts/_build.sh` to `build.sh` and change the environment variable settings inside.

### About Database

The database uses PostgreSQL, and you can start only the database using `docker-compose.dev.yaml`.

### Installing Go Dependencies

```bash
go mod tidy
```

### Running Only

```bash
cd src
../scripts/build.sh --dev
```

### Building

```bash
cd src
../scripts/build.sh
```

## Built With

- [disgo](https://pkg.go.dev/github.com/disgoorg/disgo) - The Discord API wrapper for Golang.
- [ohraban/opus](https://pkg.go.dev/github.com/hraban/opus) - The Golang bindings for the xiph.org C libraries libopus and libopusfile.
- [gorm](https://gorm.io/ja_JP/) - The fantastic ORM library for Golang.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

For the versions available, see the [tags on this repository](https://github.com/yuito-it/UntitledBot/tags).

## Authors

- @yuito-it

See also the list of [contributors](https://github.com/unipro-tech/unibot/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

- Hat tip to anyone whose code was used
- Inspiration
- etc
