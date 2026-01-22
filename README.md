# UniBot

A Discord Bot that manages and operates within All-Japan Digital Creative Club UniProject.

## Environment Variables

- `DISCORD_TOKEN` - Discord Token
- `CONFIG_ADMIN_GUILD_ID` - Admin Guild ID
- `CONFIG_ADMIN_ROLE_ID` - Admin Role ID
- `PG_DSN` - DSN for PostgreSQL connection
- `VOICEVOX_URI` - (Optional) VOICEVOX Engine URI. Default: localhost:53000
- `VOICEVOX_API_KEY` - (Optional) API Key passed as Authorization: ApiKey when communicating with VOICEVOX Engine

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

> [!TIPS]
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

- [discordgo](https://pkg.go.dev/github.com/bwmarrin/discordgo) - The Discord SDK for Golang.
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
