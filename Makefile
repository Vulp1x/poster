LOCAL_BIN:=$(CURDIR)/bin
GO_VERSION:=$(shell go version)
GO_VERSION_SHORT:=$(shell echo $(GO_VERSION) | sed -E 's/.* go(.*) .*/\1/g')

SERVICE_PATH=github.com/Vulp1x/instabot
export GO111MODULE=on
export GOPROXY=direct

DB_DSN:=$(shell grep 'pg-dsn' deploy/configs/values_local.yaml | sed "s/.*pg-dsn: //g" |  sed "s/\"//g")
PSQL_DSN:=$(shell echo $(DB_DSN) | sed 's/&timezone=utc//g')
ROOT_PSQL_DSN:=$(shell echo $(PSQL_DSN) | sed 's/insta_poster/postgres/g')

##################### GOLANG-CI RELATED CHECKS #####################
# Check global GOLANGCI-LINT
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG:=1.44.2

# Check local bin version
ifneq ($(wildcard $(GOLANGCI_BIN)),)
GOLANGCI_BIN_VERSION:=$(shell $(GOLANGCI_BIN) --version)
ifneq ($(GOLANGCI_BIN_VERSION),)
GOLANGCI_BIN_VERSION_SHORT:=$(shell echo "$(GOLANGCI_BIN_VERSION)" | sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_BIN_VERSION_SHORT:=0
endif
ifneq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_BIN_VERSION_SHORT)))"
GOLANGCI_BIN:=$(shell which golangci-lint)
endif
endif

##################### GOLANG-CI RELATED CHECKS #####################

#####################    GO RELATED CHECKS     #####################
# We always use go 1.17+
ifneq ("1.17","$(shell printf "$(GO_VERSION_SHORT)\n1.17" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.17. Found: $(GO_VERSION_SHORT))
endif
#####################    GO RELATED CHECKS     #####################

# install golangci-lint binary
.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint v$(GOLANGCI_TAG))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif



bin-deps:
	GOBIN=$(LOCAL_BIN) go install goa.design/goa/v3/cmd/goa@v3 && \
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.2.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.15.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/goresed/goresed/cmd/goresed@v0.2.3 && \
	GOBIN=$(LOCAL_BIN) go install github.com/mikefarah/yq/v4@v4.22.1



.PHONY: generate-db
generate-db: migrate-db generate-db-structure generate-db-code


generate-db-structure:
	pg_dump "$(PSQL_DSN)" --schema-only --no-owner --no-privileges --no-tablespaces --no-security-labels --exclude-table=goose_db_version > structure.sql


generate-db-code:
	$(LOCAL_BIN)/sqlc generate --file=sqlc.yaml
	$(LOCAL_BIN)/goresed regenerate --file=goresed/sqlc.yaml --references=goresed/sqlc_entry_point_ref.yaml

migrate-db:
	ENVIRONMENT=local scripts/migrate.sh



regenerate-db-drop-create:
	psql "$(ROOT_PSQL_DSN)" --command="DROP DATABASE IF EXISTS insta_poster;"
	psql "$(ROOT_PSQL_DSN)" --command="CREATE DATABASE insta_poster WITH OWNER postgres ENCODING = 'UTF8';"


regenerate-db: regenerate-db-drop-create migrate-db generate-db-structure generate-db-code

# run diff lint like in pipeline
.PHONY: lint
lint: install-lint
	$(info Running lint...)
	$(GOLANGCI_BIN) run --new-from-rev=origin/master --config=.golangci.yml ./...


build:
	go mod download && CGO_ENABLED=0  go build \
		-o ./bin/rest-api ./cmd/rest-api/main.go

#generate code from design.go
.PHONY: gen
gen:
	$(LOCAL_BIN)/goa gen github.com/inst-api/poster/design
	#$(LOCAL_BIN)/goa example github.com/inst-api/poster/design -o internal/service
	$(LOCAL_BIN)/yq -i '.servers[0].url = "/"'  gen/http/openapi3.yaml
	$(LOCAL_BIN)/yq -i '.host = "/"' gen/http/openapi.yaml
	rm -rf ./internal/service/cmd
	rm -rf ./internal/service/gen
	rm -rf ./internal/service/cli
	rm -rf ./gen/http/*/client
	rm -rf ./gen/*/client.go
	rm -rf ./gen/http/cli
	git restore ./gen/tasks_service/consts.go


.PHONY: migrate-prod
migrate-prod:
	 ENVIRONMENT=production scripts/migrate.sh

.PHONY: migrate-local
migrate-local:
	 ENVIRONMENT=local scripts/migrate.sh


.PHONY: migrate-test
migrate-test:
	 ENVIRONMENT=test scripts/migrate.sh

restart-test:
	git pull && \
	docker-compose down && \
	docker-compose -f docker-compose-local.yml up --build --detach && \
	sleep 30 &&\
	make migrate-test

hard-restart-test: drop-pgdata restart-test

drop-pgdata:
	rm -rf ./pgdata