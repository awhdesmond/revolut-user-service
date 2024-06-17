# Revolut User Service

This service provides a HTTP-based API for managing users' date of birth.

## Dependencies

| Dependency | Version |
| ---------- | ------- |
| Go         | 1.20    |
| Postgres   | 16.2    |
| Redis      | 7.0     |

## Getting Started

```bash
# Deploy Postgres, Redis, Swagger OpenAPI
# and HTTP API using docker-compose
docker-compose up -d

# Perform SQL migrations on postgres
brew install flyway
make db

# Run simple queries against HTTP API
./scripts/simple-query.sh

# Rebuild container image when you made changes
docker-compose up -d --no-deps --build api
```

To run binary on local machine:

```bash
make build

# or use direnv (https://direnv.net/)
cp .envrc.template .envrc; export $(cat .envrc | xargs)
./build/server
```

## Testing
```bash
make test-db
make test
```

## Docker
```bash
export DOCKER_REPOSITORY=<REPO_URL>
make docker
make docker-push
```

## OpenAPI

View the OpenAPI spec for this service at http://localhost:3000.

## GitHub Actions

View the GitHub Actions Workflows under `.github` directory.
