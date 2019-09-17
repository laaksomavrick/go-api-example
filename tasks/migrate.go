package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
)

// TODO dockerize so no dependency on psql

func main() {

	// TODO os.env()
	log.Print("Running migrations from migrations")

	db, err := sqlx.Connect("postgres", "user=postgres sslmode=disable")

	if err != nil {
		panic(err)
	}

	filenames := getMigrationFilenames()

	for _, filename := range filenames {
		log.Printf("Found %s", filename)
		fileContent := getMigrationFileContent(filename)
		db.MustExec(fileContent)
	}

	log.Print("Migrations complete")

}

func getMigrationFilenames() []string {
	// TODO os.env()
	var filenames []string
	files, err := ioutil.ReadDir("migrations")

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames
}

func getMigrationFileContent(filename string) string {
	// TODO os.env()
	path := fmt.Sprintf("migrations/%s", filename)
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(data)
}