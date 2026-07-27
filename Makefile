.PHONY: migration-apply migration-diff schema-apply schema-lint

migration-apply:
	set -a && source .env && set +a && atlas migrate apply --env "local"
migration-diff:
	set -a && source .env && set +a && atlas migrate diff --env "local"
migration-hash:
	set -a && source .env && set +a && atlas migrate hash
schema-apply:
	set -a && source .env && set +a && atlas schema apply --env "local"
schema-lint:
	set -a && source .env && set +a && atlas schema lint --env "local"
db-clean:
	set -a && source .env && set +a && atlas schema clean --env "local"

model-gen:
	set -a && source .env && set +a && cd src && go run cmd/gen/main.go && cd ../ 

COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git branch --show-current)
VERSION := $(shell git describe --tags --abbrev=0)

LDFLAGS = \
	-X unibot/internal.GitCommit=$(COMMIT) \
	-X unibot/internal.GitBranch=$(BRANCH)

BUILD_LDFLAGS = \
	-X unibot/internal.Version=$(VERSION) \
	$(LDFLAGS)

BOT_TARGET = cmd/bot/main.go
RSS_TARGET = cmd/rss_cron/main.go

.PHONY: run run-rss build build-rss

run:
	cd ./src && go run -ldflags "$(LDFLAGS)" $(BOT_TARGET)

run-rss:
	cd ./src && go run -ldflags "$(LDFLAGS)" $(RSS_TARGET)

build:
	cd ./src && go build -ldflags "$(BUILD_LDFLAGS)" $(BOT_TARGET)

build-rss:
	cd ./src && go build -ldflags "$(BUILD_LDFLAGS)" $(RSS_TARGET)