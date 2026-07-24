.PHONY: migration-apply migration-diff schema-apply schema-lint

migration-apply:
	atlas migrate apply --env "local"
migration-diff:
	atlas migrate diff --env "local"
migration-hash:
	atlas migrate hash
schema-apply:
	atlas schema apply --env "local"
schema-lint:
	atlas schema lint --env "local"
db-clean:
	atlas schema clean --env "local"
	atlas exec --env "local" "CREATE SCHEMA IF NOT EXISTS public;"
