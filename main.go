package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"time"

	"alecgibson.ca/go-postgres-petstore/pkg/infrastructure/db"
	"alecgibson.ca/go-postgres-petstore/pkg/interface/api"
	"alecgibson.ca/go-postgres-petstore/pkg/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

const (
	timeout                 = 15 * time.Second
	defaultConnectionString = "postgresql://postgres:password@localhost:5432"
	defaultMigrateDB        = "postgres://postgres:password@localhost:5432?sslmode=disable"
	connectionStringKey     = "CONNECTION_STRING"
	databaseStringKey       = "MIGRATE_DATABASE"
)

var (
	//go:embed sql/migrations/*
	embeddedFS embed.FS
)

func main() {
	doDatabaseMigrations()

	connectionString, found := os.LookupEnv(connectionStringKey)
	if !found {
		connectionString = defaultConnectionString
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connection, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	querier := db.New(connection)
	petService := service.NewPetService(querier)
	apiController := api.NewController(petService)

	e := echo.New()
	api.RegisterHandlers(e, apiController)
	e.Logger.Fatal(e.Start(":5000"))
}

func doDatabaseMigrations() {
	migrateDB, found := os.LookupEnv(databaseStringKey)
	if !found {
		migrateDB = defaultMigrateDB
	}

	source, err := iofs.New(embeddedFS, "sql/migrations")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to setup embedded filesystem source")
	}

	migrator, err := migrate.NewWithSourceInstance("embed", source, migrateDB)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to initialize migrate")
	}

	err = migrator.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("All schema migrations applied")
		} else {
			panic("Failed to apply schema migrations")
		}
	}
}
