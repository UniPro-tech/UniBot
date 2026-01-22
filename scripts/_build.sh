COMMIT=$(git rev-parse --short HEAD)
BRANCH=$(git branch --show-current)

if [[ "$*" == *"--dev"* ]]; then
  DISCORD_TOKEN="your_token_here"
  CONFIG_ADMIN_GUILD_ID="your_guild_id_here"
  CONFIG_ADMIN_ROLE_ID="your_role_id_here"
  PG_DSN="your_postgres_dsn_here"

  go run src/cmd/bot/main.go -ldflags "\
-X unibot/internal.GitCommit=$COMMIT \
-X unibot/internal.GitBranch=$BRANCH" \
cmd/bot/main.go
else
  VERSION=$(git describe --tags --abbrev=0)

  go build -ldflags "\
-X unibot/internal.Version=$VERSION \
-X unibot/internal.GitCommit=$COMMIT \
-X unibot/internal.GitBranch=$BRANCH" \
cmd/bot/main.go
fi