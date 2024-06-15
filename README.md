# Revolut Users Service

This service prvovdes a HTTP-based API for managing users birthday.

## Dependencies

## Getting Started

```bash
# Use docker-compose to deploy postgres, redis and api server
docker-compose up

# Perform SQL migrations on postgres
brew install flyway
make db
```

## Testing
```bash
make test
```
