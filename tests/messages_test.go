package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-palindrome/tests/helpers"
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
}

// BeforeEach
func (suite *MessageTestSuite) SetupTest() {
	err := suite.Truncate("messages")
	if err != nil {
		suite.Fail(err.Error())
	}
}

func (suite *MessageTestSuite) TestSomething() {
	assert.Equal(suite.T(), 200, 200)
}

// healthz
// 200s

// get all messages
// 200s
// sends empty list when db is empty
// sends n records when db has n records

// create message
// 400s for a bad request (e.g something not json)
// 422s for content length 0, length 1025
// 200s for message; has proper shape

// get message
// 400s for bad url
// 404s for message that doesn't exist
// 200s for message that does exist; message has proper shape

// update message
// 400s for bad url
// 400s for non-json
// 422s for content 0, 1025
// 404s for not found
// 200s for ok; message has proper shape; updatedAt != createdAt

// delete message
// 400s for bad url
// 404s for doesn't exist
// 200s for ok
