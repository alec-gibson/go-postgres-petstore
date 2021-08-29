//go:generate mockgen -destination=mock_repository.go -package=db alecgibson.ca/go-postgres-petstore/pkg/infrastructure/db Querier
package db
