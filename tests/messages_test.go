package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-palindrome/tests/helpers"
	"io"
	"net/http"
	"testing"
)

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}

type MessageTestSuite struct {
	helpers.IntegrationTestSuite
}

// BeforeAll
func (suite *MessageTestSuite) SetupSuite() {
	suite.Init()
	suite.CheckTableExists("messages")
}

// AfterAll
func (suite *MessageTestSuite) TearDownSuite() {
	err := suite.Truncate("messages")
	if err != nil {
		suite.Fail(err.Error())
	}
}

// BeforeEach
func (suite *MessageTestSuite) SetupTest() {
	err := suite.Truncate("messages")
	if err != nil {
		suite.Fail(err.Error())
	}
}

// GET /healthz

func (suite *MessageTestSuite) TestHealthzOk() {
	url := suite.GetApiUrl("healthz")

	response, err := http.Get(url)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), "ok", body["resource"])
	assert.Equal(suite.T(), nil, body["error"])
}

// GET /messages

func (suite *MessageTestSuite) TestGetAllMessagesOkEmpty() {
	url := suite.GetApiUrl("messages")

	response, err := http.Get(url)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), []interface{}{}, body["resource"])
	assert.Equal(suite.T(), nil, body["error"])
}

func (suite *MessageTestSuite) TestGetAllMessagesOkWithMessages() {
	url := suite.GetApiUrl("messages")

	response, err := http.Get(url)

	// Insert some data to confirm resource actually returns all messages
	_ = suite.seedMessage("foo", false)
	_ = suite.seedMessage("foooof", true)

	if err != nil {
		suite.HandleError(err)
	}

	response, err = http.Get(url)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	messages := body["resource"].([]interface{})

	firstMessage := messages[0].(map[string]interface{})
	secondMessage := messages[1].(map[string]interface{})

	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), nil, body["error"])

	assert.Equal(suite.T(), 2, len(messages))

	assert.Equal(suite.T(), "foo", firstMessage["content"])
	assert.Equal(suite.T(), false, firstMessage["isPalindrome"])

	assert.Equal(suite.T(), "foooof", secondMessage["content"])
	assert.Equal(suite.T(), true, secondMessage["isPalindrome"])
}

// POST /message

func (suite *MessageTestSuite) TestCreateMessageBadRequest() {
	url := suite.GetApiUrl("message")

	// nil is sent, endpoint is expecting json payload
	response, err := http.Post(url, "application/json", nil)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 400, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestCreateMessageUnprocessableEntityRequestTooShort() {
	url := suite.GetApiUrl("message")

	// Content length must be > 0 and <= 1024
	requestBody := map[string]interface{}{
		"content": "",
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	// nil is sent, endpoint is expecting json payload
	response, err := http.Post(url, "application/json", requestBodyBuffer)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 422, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestCreateMessageUnprocessableEntityRequestTooLong() {
	longString := suite.generateLongString()
	url := suite.GetApiUrl("message")

	// Content length must be > 0 and <= 1024
	requestBody := map[string]interface{}{
		"content": longString,
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	// nil is sent, endpoint is expecting json payload
	response, err := http.Post(url, "application/json", requestBodyBuffer)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 422, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestCreateMessageOkIsPalindrome() {
	url := suite.GetApiUrl("message")

	requestBody := map[string]interface{}{
		"content": "helloolleh",
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	response, err := http.Post(url, "application/json", requestBodyBuffer)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode


	message := body["resource"].(map[string]interface{})

	assert.Equal(suite.T(), 200, statusCode)
	assert.Nil(suite.T(), body["error"])

	assert.NotNil(suite.T(), message["id"])
	assert.NotNil(suite.T(), message["createdAt"])
	assert.NotNil(suite.T(), message["updatedAt"])

	assert.Equal(suite.T(), true, message["isPalindrome"])
	assert.Equal(suite.T(), "helloolleh", message["content"])
}

func (suite *MessageTestSuite) TestCreateMessageOkIsNotPalindrome() {
	url := suite.GetApiUrl("message")

	requestBody := map[string]interface{}{
		"content": "foobar",
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	response, err := http.Post(url, "application/json", requestBodyBuffer)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode


	message := body["resource"].(map[string]interface{})

	assert.Equal(suite.T(), 200, statusCode)
	assert.Nil(suite.T(), body["error"])

	assert.NotNil(suite.T(), message["id"])
	assert.NotNil(suite.T(), message["createdAt"])
	assert.NotNil(suite.T(), message["updatedAt"])

	assert.Equal(suite.T(), false, message["isPalindrome"])
	assert.Equal(suite.T(), "foobar", message["content"])
}

// GET /message/{id}

func (suite *MessageTestSuite) TestGetMessageBadRequest() {
	badUrl := suite.GetApiUrl("message/asdasd")

	response, err := http.Get(badUrl)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 400, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestGetMessageNotFoundRequest() {
	badUrl := suite.GetApiUrl("message/999")

	response, err := http.Get(badUrl)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 404, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestGetMessageOk() {
	id := suite.seedMessage("foo", false)

	url := suite.GetApiUrl(fmt.Sprintf("message/%d", id))

	response, err := http.Get(url)

	if err != nil {
		suite.HandleError(err)
	}

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	message := body["resource"].(map[string]interface{})

	assert.Equal(suite.T(), 200, statusCode)
	assert.Nil(suite.T(), body["error"])

	assert.NotNil(suite.T(), message["id"])
	assert.NotNil(suite.T(), message["createdAt"])
	assert.NotNil(suite.T(), message["updatedAt"])
	assert.NotNil(suite.T(), message["isPalindrome"])
	assert.NotNil(suite.T(), message["content"])
}

// PATCH /message/{id}

func (suite *MessageTestSuite) TestUpdateMessageBadRequestUrl() {
	badUrl := suite.GetApiUrl("message/asdasd")

	response := suite.sendPatch(badUrl, nil)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 400, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestUpdateMessageBadRequestBody() {
	badUrl := suite.GetApiUrl("message/1")

	response := suite.sendPatch(badUrl, nil)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 400, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestUpdateMessageUnprocessableEntityRequestTooLong() {
	id := suite.seedMessage("foo", false)
	url := suite.GetApiUrl(fmt.Sprintf("message/%d", id))

	// Content length must be > 0 and <= 1024
	requestBody := map[string]interface{}{
		"content": suite.generateLongString(),
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	response := suite.sendPatch(url, requestBodyBuffer)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 422, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestUpdateMessageUnprocessableEntityRequestTooShort() {
	id := suite.seedMessage("foo", false)
	url := suite.GetApiUrl(fmt.Sprintf("message/%d", id))

	// Content length must be > 0 and <= 1024
	requestBody := map[string]interface{}{
		"content": "",
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	response := suite.sendPatch(url, requestBodyBuffer)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 422, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestUpdateMessageNotFoundRequest() {
	badUrl := suite.GetApiUrl("message/9999")

	requestBody := map[string]interface{}{
		"content": "something",
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	response := suite.sendPatch(badUrl, requestBodyBuffer)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 404, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestUpdateMessageOk() {
	id := suite.seedMessage("foo", false)
	url := suite.GetApiUrl(fmt.Sprintf("message/%d", id))

	requestBody := map[string]interface{}{
		"content": "helloolleh",
	}

	requestBodyBuffer := suite.MapToBuffer(requestBody)

	response := suite.sendPatch(url, requestBodyBuffer)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	message := body["resource"].(map[string]interface{})

	assert.Equal(suite.T(), 200, statusCode)
	assert.Nil(suite.T(), body["error"])

	assert.NotNil(suite.T(), message["id"])
	assert.NotNil(suite.T(), message["createdAt"])
	assert.NotNil(suite.T(), message["isPalindrome"])

	assert.Equal(suite.T(), "helloolleh", message["content"])
	assert.Equal(suite.T(), true, message["isPalindrome"])
	assert.NotEqual(suite.T(), message["createdAt"], message["updatedAt"])
}

// DELETE /message/{id}

func (suite *MessageTestSuite) TestDeleteMessageBadRequestUrl() {
	badUrl := suite.GetApiUrl("message/qweqwe")

	response := suite.sendDelete(badUrl)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 400, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestDeleteMessageNotFoundRequest() {
	badUrl := suite.GetApiUrl("message/1337")

	response := suite.sendDelete(badUrl)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 404, statusCode)
	assert.NotNil(suite.T(), body["error"])
}

func (suite *MessageTestSuite) TestDeleteMessageOk() {
	id := suite.seedMessage("foo", false)
	url := suite.GetApiUrl(fmt.Sprintf("message/%d", id))

	response := suite.sendDelete(url)

	body := suite.ResponseToMap(response)
	statusCode := response.StatusCode

	assert.Equal(suite.T(), 200, statusCode)
	assert.Nil(suite.T(), body["error"])
}

// Helpers

func (suite *MessageTestSuite) seedMessage(content string, isPalindrome bool) int {
	var id int

	err := suite.DB.QueryRow("INSERT INTO messages (content, is_palindrome) VALUES ($1, $2) RETURNING id", content, isPalindrome).Scan(&id)

	if err != nil {
		suite.HandleError(err)
	}

	return id
}

func (suite *MessageTestSuite) sendPatch(url string, body io.Reader) *http.Response {
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPatch, url, body)

	if err != nil {
		suite.HandleError(err)
	}

	response, err := client.Do(request)

	if err != nil {
		suite.HandleError(err)
	}

	return response
}

func (suite *MessageTestSuite) sendDelete(url string) *http.Response {
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		suite.HandleError(err)
	}

	response, err := client.Do(request)

	if err != nil {
		suite.HandleError(err)
	}

	return response
}

func (suite *MessageTestSuite) generateLongString() string {
	longString := ""

	for i := 0; i <= 1024; i++ {
		longString += "a"
	}

	return longString
}
