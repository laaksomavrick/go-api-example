package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"go-palindrome/lib"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	driver = "postgres"
)

// IntegrationTestSuite defines the shape of an integration test suite, providing common
// functionality (connecting to the database) alongside some utility functions required for
// integration testing.
type IntegrationTestSuite struct {
	suite.Suite
	DB           *sqlx.DB
	apiHost      string
	apiPort      string
	postgresUser string
	postgresHost string
}

// Init initializes the test suite, establishing a database connection and setting up config values.
func (suite *IntegrationTestSuite) Init() {
	suite.apiHost = os.Getenv("API_HOST")
	suite.apiPort = os.Getenv("API_PORT")
	suite.postgresUser = os.Getenv("POSTGRES_USER")
	suite.postgresHost = os.Getenv("POSTGRES_HOST")

	connectionString := lib.GetConnectionString(suite.postgresUser, suite.postgresHost)
	db, err := lib.ConnectToDb(driver, connectionString)

	if err != nil {
		suite.HandleError(err)
	}

	suite.DB = db
}

// Truncate truncates a specified table.
func (suite *IntegrationTestSuite) Truncate(table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s", table)
	_, err := suite.DB.Exec(query)
	return err
}

// CheckTableExists checks whether a table exists, with retry logic.
func (suite *IntegrationTestSuite) CheckTableExists(table string) {
	exists := false
	tries := 0

	for {
		if tries > 10 {
			suite.Fail("unable to find table")
		}

		err := suite.DB.QueryRow(`
			SELECT EXISTS (
				SELECT * FROM information_schema.tables
				WHERE table_name = $1
			);
		`, table).Scan(&exists)

		if err != nil {
			suite.HandleError(err)
		}

		if exists == true {
			break
		}

		log.Print("sleeping while waiting for table to exist")
		tries += 1
		time.Sleep(5 * time.Second)
	}

}

// MapToBuffer maps a go map to a buffer.
func (suite *IntegrationTestSuite) MapToBuffer(body map[string]interface{}) *bytes.Buffer {
	marshalled, err := json.Marshal(body)

	if err != nil {
		suite.HandleError(err)
	}

	return bytes.NewBuffer(marshalled)
}

// ResponseToMap maps a http response to a go map.
func (suite *IntegrationTestSuite) ResponseToMap(response *http.Response) map[string]interface{} {
	var body map[string]interface{}

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		suite.HandleError(err)
	}

	err = json.Unmarshal(responseBody, &body)

	if err != nil {
		suite.HandleError(err)
	}

	return body
}

// GetApiUrl returns the api url for http requests.
func (suite *IntegrationTestSuite) GetApiUrl(path string) string {
	return fmt.Sprintf("http://%s:%s/%s", suite.apiHost, suite.apiPort, path)
}

// HandleError is a helper for handling an error in the test suite.
func (suite *IntegrationTestSuite) HandleError(err error) {
	suite.Fail(err.Error())
}
