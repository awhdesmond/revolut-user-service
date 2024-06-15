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
	./scrtips/clean-db.sh
	./scrtips/migrate-db.sh

test-db:
	docer exec postgres \
		psql -U postgres -c 'CREATE DATABASE postgres-test WITH OWNER postgres'

