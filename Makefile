SERVICES ?=
.PHONY: migrations migrations-reset migrations-fresh db-fresh seeds admin schema-sync sqlc-generate sqlc-gen sqlc-vet sqlc-verify proto-tools proto-gen api-gen frontend-api-gen test lint check backend-test frontend-test backend-lint frontend-lint


	
#
# Database: migrations and sqlc
#

# Run migrations
migrations:
	@./scripts/migrations.sh up "$(SERVICES)"

# Roll back migrations
migrations-reset:
	@./scripts/migrations.sh reset "$(SERVICES)"

# Add test data to the database
seeds:
	@./scripts/seeds.sh "$(SERVICES)"

# Create/refresh the admin user through the CLI
admin:
	@./scripts/admin.sh

# Generate code, roll back migrations, run them again, and add test data
db-fresh: sqlc-gen migrations-reset migrations seeds

# Backward-compatible alias kept for callers using the old target name.
migrations-fresh: db-fresh

schema-sync:
	@./scripts/schema-sync.sh "$(SERVICES)"

sqlc-gen:
	@./scripts/schema-sync.sh "$(SERVICES)"
	@./scripts/sqlc.sh gen "$(SERVICES)"
	
sqlc-query-gen:
	@./scripts/query-gen/crud-query-gen.sh "$(SERVICES)"

# Backward-compatible alias kept for callers using the old target name.
sqlc-generate: sqlc-gen

sqlc-vet:
	@./scripts/sqlc.sh vet "$(SERVICES)"

sqlc-verify:
	@./scripts/sqlc.sh verify "$(SERVICES)"

#
# OpenAPI / Proto generation
#
api-gen: proto-gen frontend-api-gen

proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.71.0
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.18.1
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.29.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.29.0

proto-gen: proto-tools
	PATH="$(shell go env GOPATH)/bin:$$PATH" buf dep update ./services/identity/api
	PATH="$(shell go env GOPATH)/bin:$$PATH" buf dep update ./services/strategy-registry/api
	go generate ./services/identity/api
	go generate ./services/strategy-registry/api

frontend-api-gen:
	@./scripts/frontend-api-gen.sh


# Tests

backend-test:
	@go test $$(go list ./... | grep -v '/node_modules/')

frontend-test:
	@if [ ! -f apps/web/package.json ]; then \
		echo "apps/web/package.json not found, skipping frontend tests"; \
		exit 0; \
	fi; \
	if node -e "const pkg=require('./apps/web/package.json'); process.exit(pkg.scripts && pkg.scripts.test ? 0 : 1)"; then \
		if [ ! -f package.json ]; then \
			echo "root package.json not found, cannot run workspace frontend tests"; \
			exit 1; \
		fi; \
		if [ ! -d node_modules ]; then \
			echo "root node_modules not found, run npm install at the repository root first"; \
			exit 1; \
		fi; \
		npm run test -w @stratflow/web; \
	else \
		echo "frontend test script is not configured yet, skipping frontend tests"; \
	fi

test: backend-test frontend-test

user:
	@go test ./services/strategy-registry/tests/e2e/grpc/strategies -run 'TestGRPCStrategiesCRUDGraphUser$$'

backend-lint:
	@status=0; \
	golangci-lint run ./... || status=$$?; \
	if [ -s golangci-lint-report.json ]; then \
		jq . golangci-lint-report.json > golangci-lint-report.tmp.json && mv golangci-lint-report.tmp.json golangci-lint-report.json; \
	fi; \
	exit $$status

frontend-lint:
	@if [ ! -f apps/web/package.json ]; then \
		echo "apps/web/package.json not found, skipping frontend lint"; \
		exit 0; \
	fi; \
	if node -e "const pkg=require('./apps/web/package.json'); process.exit(pkg.scripts && pkg.scripts.lint ? 0 : 1)"; then \
		if [ ! -f package.json ]; then \
			echo "root package.json not found, cannot run workspace frontend lint"; \
			exit 1; \
		fi; \
		if [ ! -d node_modules ]; then \
			echo "root node_modules not found, run npm install at the repository root first"; \
			exit 1; \
		fi; \
		npm run lint -w @stratflow/web; \
	else \
		echo "frontend lint script is not configured yet, skipping frontend lint"; \
	fi

lint: backend-lint frontend-lint

check: lint test
