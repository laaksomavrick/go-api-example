package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	log.Print("Running migrations from migrations")

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresHost := os.Getenv("POSTGRES_HOST")
	migrationsDirectory := os.Getenv("MIGRATIONS_DIRECTORY")
	driver := "postgres"

	connectionString := fmt.Sprintf("user=%s sslmode=disable host=%s", postgresUser, postgresHost)

	log.Printf("Connection string: %s", connectionString)
	log.Printf("Migrations directory: %s", migrationsDirectory)

	var db *sqlx.DB
	tries := 0

	for {

		if tries > 10 {
			log.Fatal("Unable to connect to postgres, exiting.")
		}

		var err error
		db, err = sqlx.Connect(driver, connectionString)

		if err == nil {
			break
		}

		log.Print("Sleeping while waiting on postgres")
		tries += 1
		time.Sleep(5 * time.Second)
	}

	// Grab the filenames in our migration directory
	filenames := getMigrationFilenames(migrationsDirectory)

	// Run them
	for _, filename := range filenames {
		log.Printf("Found %s", filename)
		fileContent := getMigrationFileContent(migrationsDirectory, filename)
		db.MustExec(fileContent)
	}

	log.Print("Migrations complete")

	// Exit with a happy status code
	os.Exit(0)
}

func getMigrationFilenames(migrationsDirectory string) []string {
	var filenames []string
	files, err := ioutil.ReadDir(migrationsDirectory)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames
}

func getMigrationFileContent(migrationsDirectory string, filename string) string {
	path := fmt.Sprintf("%s/%s", migrationsDirectory, filename)
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(data)
}