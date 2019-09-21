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
	"net/http"
	"os"
)

const (
	driver = "postgres"
)

type IntegrationTestSuite struct {
	suite.Suite
	DB           *sqlx.DB
	apiHost      string
	apiPort      string
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
		suite.HandleError(err)
	}

	suite.DB = db
}

func (suite *IntegrationTestSuite) Truncate(table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s", table)
	_, err := suite.DB.Exec(query)
	return err
}

func (suite *IntegrationTestSuite) MapToBuffer(body map[string]interface{}) *bytes.Buffer {
	marshalled, err := json.Marshal(body)

	if err != nil {
		suite.HandleError(err)
	}

	return bytes.NewBuffer(marshalled)
}

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

func (suite *IntegrationTestSuite) GetApiUrl(path string) string {
	return fmt.Sprintf("http://%s:%s/%s", suite.apiHost, suite.apiPort, path)
}

func (suite *IntegrationTestSuite) HandleError(err error) {
	suite.Fail(err.Error())
}
