package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-palindrome/src/palindrome"
	"log"
	_ "github.com/lib/pq"
)

const (
	driver = "postgres"
)

func main()  {
	config := palindrome.NewConfig()
	router := mux.NewRouter()

	connectionString := fmt.Sprintf("user=%s sslmode=disable host=%s", config.PostgresUser, config.PostgresHost)
	db, err := palindrome.ConnectToDb(driver, connectionString)

	if err != nil {
		log.Fatal(err)
	}

	server := palindrome.NewServer(router, db, config)

	server.Init()
}

