# Revolut User Service

This service prvovdes a HTTP-based API for managing users birthday.

## Dependencies

## Getting Started

```bash
# Use docker-compose to deploy postgres, redis and api server, and OpenAPI
mkdir -p ./.data/postgres
docker-compose up

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
make test
```
