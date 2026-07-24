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