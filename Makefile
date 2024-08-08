include .env
export

COMPOSE=docker compose -f ./docker/compose.yaml

MIGRATE=docker compose -f ./docker/compose.yaml exec web sql-migrate
ifeq ($(p),host)
	MIGRATE=sql-migrate
endif

db:
	$(MIGRATE) $(ARGS)

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
# activate-venv:
# 	source .venv/bin/activate

lint-setup: # must activate python virtual environment
	python -m ensurepip --upgrade
	pip install pre-commit
	pre-commit install
	pre-commit autoupdate

.PHONY: db migrate-status migrate-up migrate-down redo create lint-setup

# *
# docker compose
# *

dc:
	${COMPOSE} $(ARGS)

dc-up:
	${COMPOSE} up -d $(ARGS)

dc-down:
	${COMPOSE} down $(ARGS)