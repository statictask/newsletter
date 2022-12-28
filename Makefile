# Makefile for Go projects.
#
# This Makefile makes an effort to provide standard make targets, as described
# by https://www.gnu.org/prep/standards/html_node/Standard-Targets.html.
SHELL := /bin/sh

include ./rules/Makefile.*

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

################################################################################
## Standard make targets
################################################################################

.DEFAULT_GOAL := all
.PHONY: all
all: fix install

.PHONY: install
install: go-install

.PHONY: uninstall
uninstall:
	@echo "Uninstalling ${SETTINGS_PROJECT_NAME}"
	$rm -f $(GOPATH)/bin/${SETTINGS_PROJECT_NAME}

.PHONY: clean
clean: go-clean
	@echo "Deleting coverage"
	@rm -f coverage.out

.PHONY: check
check: test

.PHONY: run
run: compose-run

################################################################################
## Go-like targets
################################################################################

.PHONY: build
build: go-build

.PHONY: test
test: go-test

.PHONY: cover
cover: go-cover/text

.PHONY: cover/html
cover/html: go-cover/html

.PHONY: cover/text
cover/text: go-cover/text

################################################################################
## Linters and formatters
################################################################################

.PHONY: fix
fix: go-fix

.PHONY: lint
lint: go-lint

################################################################################
## Migrations
################################################################################

.PHONY: migrate
migrate:
	@echo ">  Migrating: UP"
	$(GORUN) -tags migrate cmd/migrate/main.go -cmd up

.PHONY: drop
drop:
	@echo ">  Migrating: DOWN"
	$(GORUN) -tags migrate cmd/migrate/main.go -cmd down

.PHONY: connect
connect:
	@echo "> Connecting to local DB"
	psql postgres://newsletter:newsletter@localhost:5432/newsletter

.PHONY: reset
reset:
	@echo ">Resetting everything"
	sudo docker-compose down --volumes
