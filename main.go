package main

import (
	"database/sql"
	"log"

	"github.com/GabrielEger2/gobanq/api"
	db "github.com/GabrielEger2/gobanq/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:postgres@db:5432/gobanq?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(testDB)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
