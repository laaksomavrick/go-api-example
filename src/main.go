package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go-palindrome/lib"
	"go-palindrome/src/palindrome"
	"log"
)

const (
	driver = "postgres"
)

func main()  {
	config := palindrome.NewConfig()
	router := mux.NewRouter()

	connectionString := lib.GetConnectionString(config.PostgresUser, config.PostgresHost)
	db, err := lib.ConnectToDb(driver, connectionString)

	if err != nil {
		log.Fatal(err)
	}

	server := palindrome.NewServer(router, db, config)

	server.Init()
}

