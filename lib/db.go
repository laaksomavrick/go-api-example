package lib

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

// GetConnectionString returns a valid connection string, used to connect to postgres.
func GetConnectionString(postgresUser string, postgresHost string) string {
	return fmt.Sprintf("user=%s sslmode=disable host=%s", postgresUser, postgresHost)
}

// ConnectToDb implements retry logic to establish a connection to postgres, returning either
// an api on top of the connection or an error. It makes 10 attempts, with a delay of 5s between
// each. Handling failure is left to the callee.
func ConnectToDb(driver string, connectionString string) (*sqlx.DB, error) {
	var db *sqlx.DB
	tries := 0

	for {
		if tries > 10 {
			return nil, errors.New("unable to connect to postgres, exiting")
		}

		var err error
		db, err = sqlx.Connect(driver, connectionString)

		if err == nil {
			break
		}

		log.Print(err)

		log.Print("sleeping while waiting on postgres")
		tries += 1
		time.Sleep(5 * time.Second)
	}

	return db, nil
}
