package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func main()  {
	fmt.Println("Hello, world")

	db, err := sqlx.Connect("postgres", "user=postgres dbname=bar sslmode=disable")
}
