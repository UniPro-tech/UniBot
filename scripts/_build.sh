COMMIT=$(git rev-parse --short HEAD)
BRANCH=$(git branch --show-current)

if [[ "$*" == *"--rss-cron"* ]]; then
  TARGET="cmd/rss_cron/main.go"
else
  TARGET="cmd/bot/main.go"
fi

if [[ "$*" == *"--dev"* ]]; then
  export DISCORD_TOKEN="your_token_here"
  export CONFIG_ADMIN_GUILD_ID="your_guild_id_here"
  export CONFIG_ADMIN_ROLE_ID="your_role_id_here"
  export PG_DSN="your_postgres_dsn_here"
  export GITHUB_OAUTH_ID="your-client-id-here"
  export GITHUB_OAUTH_SECRET="your-client-secret-here"

  go run -ldflags "\
-X unibot/internal.GitCommit=$COMMIT \
-X unibot/internal.GitBranch=$BRANCH" \
"$TARGET"
else
  VERSION=$(git describe --tags --abbrev=0)

  go build -ldflags "\
-X unibot/internal.Version=$VERSION \
-X unibot/internal.GitCommit=$COMMIT \
-X unibot/internal.GitBranch=$BRANCH" \
"$TARGET"
fi