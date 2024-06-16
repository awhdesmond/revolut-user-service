# Revolut User Service

This service provides a HTTP-based API for managing users' date of birth.

## Dependencies

## Getting Started

```bash
# Deploy local postgres, redis, and Swagger OpenAPI
docker-compose up -d

# Perform SQL migrations on postgres
brew install flyway
make db

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
