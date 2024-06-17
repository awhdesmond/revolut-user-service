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

Or to build outside of docker:

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

## OpenAPI

View the OpenAPI spec for this service at http://localhost:3000.

## GitHub Actions

View the GitHub Actions Workflows under `.github` directory.
