package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go-palindrome/lib"
	"go-palindrome/src/api"
	"log"
)

const (
	driver = "postgres"
)

// Application entry point
func main()  {
	config := api.NewConfig()
	router := mux.NewRouter()

	connectionString := lib.GetConnectionString(config.PostgresUser, config.PostgresHost)
	db, err := lib.ConnectToDb(driver, connectionString)

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(router, db, config)

	server.Init()
}

