GITCOMMIT?=$(shell git describe --dirty --always)
CGO_ENABLED?=0
BINARY:=server
LDFLAGS:="-s -w -X github.com/awhdesmond/revolut-user-service/pkg/common.GitCommit=$(GITCOMMIT)"

DOCKER_IMAGE?=revolut-user-service_api
DOCKER_REPOSITORY?=974860574511.dkr.ecr.eu-west-1.amazonaws.com/revolut-user-service

.PHONY: build clean test db

build:
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags=$(LDFLAGS) -o build/server cmd/server/*.go

test:
	go test ./... -short -timeout 120s -race -count 1 -v

test-coverage:
	go test ./... -short -timeout 120s -race -count 1 -v -coverprofile=coverage.out
	go tool cover -html=coverage.out

db:
	./scripts/clean-db.sh postgres
	./scripts/migrate-db.sh postgres

test-db:
	docker exec revolut-user-service_postgres_1 \
		psql -U postgres -c 'CREATE DATABASE postgres_test WITH OWNER postgres' || true
	./scripts/clean-db.sh postgres_test
	./scripts/migrate-db.sh postgres_test

docker:
	docker build -t $(DOCKER_IMAGE):$(GITCOMMIT) .
	docker tag $(DOCKER_IMAGE):$(GITCOMMIT) $(DOCKER_REPOSITORY):$(GITCOMMIT)

docker-push:
	docker push $(DOCKER_REPOSITORY):$(GITCOMMIT)

clean:
	rm -rf build cover.html coverage.out
