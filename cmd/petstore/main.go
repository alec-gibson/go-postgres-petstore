package main

import (
	"context"
	"os"
	"time"

	"alecgibson.ca/go-postgres-petstore/pkg/api"
	"alecgibson.ca/go-postgres-petstore/pkg/db"
	"alecgibson.ca/go-postgres-petstore/pkg/service"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

const (
	timeout                 = 15 * time.Second
	defaultConnectionString = "postgresql://postgres:password@localhost:5432"
	connectionStringKey     = "CONNECTION_STRING"
)

func main() {
	connectionString, found := os.LookupEnv(connectionStringKey)
	if !found {
		connectionString = defaultConnectionString
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connection, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		panic("Failed to connect to database")
	}

	querier := db.New(connection)
	petService := service.NewPetService(querier)
	apiController := api.NewController(petService)

	e := echo.New()
	api.RegisterHandlers(e, apiController)
	e.Logger.Fatal(e.Start(":5000"))
}
