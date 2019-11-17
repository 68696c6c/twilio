package twilio

import (
	"os"
	"strings"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_streamToByte(t *testing.T) {
	input := fake.CharactersN(36)
	r := strings.NewReader(input)
	b, err := streamToByte(r)
	require.Nil(t, err, "unexpected error returned")
	assert.Equal(t, input, string(b), "unexpected result returned")
}

func TestClient_SendSMS(t *testing.T) {
	accountSID := os.Getenv("ACCOUNT_SID")
	require.NotEmpty(t, accountSID, "missing ACCOUNT_SID")

	authToken := os.Getenv("AUTH_TOKEN")
	require.NotEmpty(t, authToken, "missing AUTH_TOKEN")

	accountPhone := os.Getenv("ACCOUNT_PHONE")
	require.NotEmpty(t, accountPhone, "missing ACCOUNT_PHONE")

	toPhone := os.Getenv("TEST_TO_PHONE")
	require.NotEmpty(t, toPhone, "missing TEST_TO_PHONE")

	c := NewClient(accountSID, authToken, accountPhone)
	response, err := c.Send(toPhone, "hello world")

	println("response: ", response.String())
	require.Nil(t, err, "unexpected error returned")
	assert.NotEqual(t, Response{}, response, "empty twilio response returned")
	assert.NotNil(t, response.HTTPResponse, "empty http response returned")
}
