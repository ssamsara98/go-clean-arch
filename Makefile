include .env
export

MIGRATE=docker compose -f ./docker/compose.yml exec web sql-migrate

ifeq ($(p),host)
	MIGRATE=sql-migrate
endif

migrate-status:
	$(MIGRATE) status

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

seed-status:
	$(MIGRATE) status --env=seed

seed-up:
	$(MIGRATE) up --env=seed

seed-down:
	$(MIGRATE) down --env=seed

redo:
	@read -p  "Are you sure to reapply the last migration? [y/n]" -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; \
	then \
		$(MIGRATE) redo; \
	fi

create:
	@read -p  "What is the name of migration?" NAME; \
	${MIGRATE} new $$NAME

create-venv:
	python3 -m venv .venv

lint-setup:
	python3 -m ensurepip --upgrade
	pip3 install pre-commit
	pre-commit install
	pre-commit autoupdate

.PHONY: migrate-status migrate-up migrate-down redo create lint-setup

# docker compose
dc:
	docker compose -f ./docker/compose.yml $(ARGS)

dc-up:
	docker compose -f ./docker/compose.yml up -d $(ARGS)

dc-down:
	docker compose -f ./docker/compose.yml down $(ARGS)
