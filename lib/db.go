package lib

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func GetConnectionString(postgresUser string, postgresHost string) string {
	return fmt.Sprintf("user=%s sslmode=disable host=%s", postgresUser, postgresHost)
}

// TODO: in production app, this could fail fast and have the pod be rescheduled
// TODO: how to handle connection being lost?
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
