package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"go-palindrome/lib"
	"io/ioutil"
	"log"
	"os"
)

// For a proper migration tool, we'd have a 'migrations' table with a record of migrations
// previously run, using that info to only run newly discovered files in the MIGRATIONS_DIRECTORY

// Since this is a demo, every execution recreates the world (ie. drops schema and recreates tables)
// for convenience

func main() {
	log.Print("running migrations")

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresHost := os.Getenv("POSTGRES_HOST")
	migrationsDirectory := os.Getenv("MIGRATIONS_DIRECTORY")
	driver := "postgres"

	connectionString := lib.GetConnectionString(postgresUser, postgresHost)

	log.Printf("connection string: %s", connectionString)
	log.Printf("migrations directory: %s", migrationsDirectory)

	// Make sure we can connect to the db
	db, err := lib.ConnectToDb(driver, connectionString)

	if err != nil {
		log.Fatal(err)
	}

	// Grab the filenames in our migration directory
	filenames, err := getMigrationFilenames(migrationsDirectory)

	if err != nil {
		log.Fatal(err)
	}

	// Run them
	for _, filename := range filenames {
		log.Printf("Found %s", filename)
		fileContent, err := getMigrationFileContent(migrationsDirectory, filename)

		if err != nil {
			log.Fatal(err)
		}

		db.MustExec(fileContent)
	}

	log.Print("migrations complete")

	// Exit with a happy status code
	os.Exit(0)
}

func getMigrationFilenames(migrationsDirectory string) ([]string, error) {
	var filenames []string
	files, err := ioutil.ReadDir(migrationsDirectory)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames, nil
}

func getMigrationFileContent(migrationsDirectory string, filename string) (string, error) {
	path := fmt.Sprintf("%s/%s", migrationsDirectory, filename)
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
