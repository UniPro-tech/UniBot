COMMIT=$(git rev-parse --short HEAD)
BRANCH=$(git branch --show-current)

if [[ "$*" == *"--dev"* ]]; then
  go build -ldflags "\
-X unibot/internal.GitCommit=$COMMIT \
-X unibot/internal.GitBranch=$BRANCH" \
cmd/bot/main.go
else
  VERSION=$(git describe --tags --abbrev=0)

  go build -ldflags "\
-X unibot/internal.Version=$VERSION \
-X unibot/internal.GitCommit=$COMMIT \
-X unibot/internal.GitBranch=$BRANCH"
fi