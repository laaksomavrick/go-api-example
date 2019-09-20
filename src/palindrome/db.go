package palindrome

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

// TODO: reconcile with connectToDb logic in migrate.go
// TODO: in production app, this could fail fast and have the pod be rescheduled
func ConnectToDb(driver string, connectionString string) (*sqlx.DB, error) {
	var db *sqlx.DB
	tries := 0

	log.Print(connectionString)

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
