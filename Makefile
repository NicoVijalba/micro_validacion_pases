SHELL := /bin/bash
APP_NAME := validacion-pases
GO ?= go

.PHONY: fmt lint test test-integration coverage run docker migrate generate-openapi vuln

fmt:
	$(GO) fmt ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint no esta instalado localmente; usando contenedor..."; \
		docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.63.4 golangci-lint run ./...; \
	fi

vuln:
	govulncheck ./...

test:
	$(GO) test -race -coverprofile=coverage.out ./internal/... ./pkg/...

test-integration:
	$(GO) test -v ./tests/integration/...

coverage: test
	$(GO) tool cover -func=coverage.out

run:
	$(GO) run ./cmd/api

docker:
	docker build -t $(APP_NAME):local .

migrate:
	docker run --rm --network host -v $(PWD)/migrations:/migrations migrate/migrate:v4.17.1 \
	  -path=/migrations -database "mysql://app:app@tcp(localhost:3306)/validacion_pases" up

generate-openapi:
	@echo "OpenAPI source of truth: docs/openapi.yaml"
