GITCOMMIT?=$(shell git describe --dirty --always)
CGO_ENABLED?=0
BINARY:=server
LDFLAGS:="-s -w -X github.com/awhdesmond/revolut-user-service/pkg/common/version.GitCommit=$(GITCOMMIT)"

.PHONY: build clean test db

build:
	go build -ldflags=$(LDFLAGS) -o build/server cmd/server/*.go

test:
	go test ./... -short -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:
	rm -rf build cover.html coverage.out

db:
	./scripts/clean-db.sh postgres
	./scripts/migrate-db.sh postgres

test-db:
	docker exec revolut-user-service_postgres_1 \
		psql -U postgres -c 'CREATE DATABASE postgres_test WITH OWNER postgres' || true
	./scripts/clean-db.sh postgres_test
	./scripts/migrate-db.sh postgres_test
