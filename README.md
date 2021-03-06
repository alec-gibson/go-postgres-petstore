# Pet Store Example Application

An implementation of the sample petstore application from the examples in [oapi-codegen](https://github.com/deepmap/oapi-codegen)

## Packages Used:
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [migrate](https://github.com/golang-migrate/migrate)

## Run the Application Locally
1. Run `docker build -t petstore .` to compile a docker image of the application
2. Run `docker-compose up -d` in the root of the repo
4. The application will now be available at localhost:5000

## Changing the Database Schema:
1. Run `./migrate.sh create <migration_name>` to create your new migration scripts
2. Fill in the new migration scripts created under `sql/migrations`
3. Regenerate go code for interacting with the database using `generate-code.sh`
4. All schema migrations will be performed automatically whenever the application starts

## Changing the API:
1. Make your updates to the openapi spec in openapi/openapi.yaml
2. Generate go code using `generate-code.sh`
3. The newly generated interfaces for API endpoint handlers will be in `pkg/api/controller.go`

## Adding new Database Queries:
1. Update `sql/queries/queries.sql` to add your new queries, or add additional files in the same directory
2. Regenerate go code for interacting with the database using `generate-code.sh`
3. The newly generated query functions will be in `pkg/db/queries.sql.go`

## Running the Integration Tests:
Run the `run-integration-tests.sh` script. This script performs the following steps:
1. It builds a new docker container for the application
2. It brings up (or recreates) all containers using docker compose
3. It waits for the application to become healthy
4. It runs all integration tests (defined in the `integrationtests` directory)

## TODO:
1. Add unit tests
2. Add pipeline logic to build / test / upload the application container on every master commit
3. Add pipeline logic to deploy the application to the cloud on every master commit
4. Expand the scope of the application, so there are multiple tables and queries with joins
5. Fork this repo and make something actually useful
