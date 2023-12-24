SHELL := /bin/bash
SCRIPT_DIRECTORY := ./script
MIGRATION_DIRECTORY := ./docs/migrations

.PHONY: migration-up
migration-up: 
	@$(SHELL) $(SCRIPT_DIRECTORY)/migration.sh "$(username)" "$(password)" "$(dbHost)" "$(databaseName)" "up" "$(qty)" "$(MIGRATION_DIRECTORY)"

.PHONY: migration-down
migration-down: 
	@$(SHELL) $(SCRIPT_DIRECTORY)/migration.sh "$(username)" "$(password)" "$(dbHost)" "$(databaseName)" "down" "$(qty)" "$(MIGRATION_DIRECTORY)"

.PHONY: migration-drop
migration-drop: 
	@$(SHELL) $(SCRIPT_DIRECTORY)/migration.sh "$(username)" "$(password)" "$(dbHost)" "$(databaseName)" "drop -f" "$(MIGRATION_DIRECTORY)"

.PHONY: migration-create
migration-create: 
	@$(SHELL) $(SCRIPT_DIRECTORY)/create_migration.sh "$(MIGRATION_DIRECTORY)" "$(migrationName)"

.PHONY: run
run: swaggo
	@go build -o build/"$(service)" -ldflags="-X main.service="$(service)""
	@./build/"$(service)"