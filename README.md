# Revolut User Service

This service provides a HTTP-based API for managing users' date of birth.

The structure of the code follows the Clean Architecture approach. In addition, we added a write-through cache since by observation, services that manages users's profile information tends to be read heavy.

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

## Environment Variables

| Environment Variable                 | Description                                           |
| ------------------------------------ | ----------------------------------------------------- |
| REVOLUT_USERS_SVC_HOST               | Host to expose HTTP API server                        |
| REVOLUT_USERS_SVC_METRICS_PORT       | Port to export HTTP API server                        |
| REVOLUT_USERS_SVC_LOG_LEVEL          | Log Level                                             |
| REVOLUT_USERS_SVC_CORS_ORIGIN        | CORS Origin                                           |
| REVOLUT_USERS_SVC_POSTGRES_HOST      | Postgres Host                                         |
| REVOLUT_USERS_SVC_POSTGRES_PORT      | Postgres Port                                         |
| REVOLUT_USERS_SVC_POSTGRES_USERNAME  | Postgres Username                                     |
| REVOLUT_USERS_SVC_POSTGRES_PASSWORD  | Postgres Password                                     |
| REVOLUT_USERS_SVC_POSTGRES_DATABASE  | Postgres Database                                     |
| REVOLUT_USERS_SVC_REDIS_URI          | Redis URI                                             |
| REVOLUT_USERS_SVC_REDIS_PASSWORD     | Redis Password                                        |
| REVOLUT_USERS_SVC_REDIS_CLUSTER_MODE | Redis Cluster Mode. Use non-empty string to enable it |


## Testing

Run the following commands to bootstrap the test database and run unit tests.

```bash
make test-db
make test
```

## Docker

Build docker image and push them to the repository.

```bash
export CONTAINER_REGISTRY=<REGISTRY_URL>
export CONTAINER_REPOSITORY=<REPOSITORY>

make docker
make docker-push
```

## Swagger OpenAPI

View the OpenAPI spec for this service at http://localhost:3000.

## GitHub Actions (CI/CD)

View the GitHub Actions Workflows (CI Pipelines) under `.github` directory.
