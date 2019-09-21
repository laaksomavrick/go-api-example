package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"go-palindrome/lib"
	"log"
	"os"
	_ "github.com/lib/pq"
)

const (
	driver = "postgres"
)

type IntegrationTestSuite struct {
	suite.Suite
	db *sqlx.DB
	apiHost string
	apiPort string
	postgresUser string
	postgresHost string
}

func (suite *IntegrationTestSuite) Init() {
	suite.apiHost = os.Getenv("API_HOST")
	suite.apiPort = os.Getenv("API_PORT")
	suite.postgresUser = os.Getenv("POSTGRES_USER")
	suite.postgresHost = os.Getenv("POSTGRES_HOST")

	connectionString := lib.GetConnectionString(suite.postgresUser, suite.postgresHost)
	db, err := lib.ConnectToDb(driver, connectionString)

	if err != nil {
		log.Print(err)
		suite.Fail("error connecting to postgres")
	}

	suite.db = db
}

func (suite *IntegrationTestSuite) Truncate(table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s", table)
	_, err := suite.db.Exec(query)
	return err
}

func (suite *IntegrationTestSuite) MapToBuffer(body map[string]string) *bytes.Buffer {
	marshalled, err := json.Marshal(body)

	if err != nil {
		log.Print(err)
		suite.Fail("error mapping body to buffer")
	}

	return bytes.NewBuffer(marshalled)
}
