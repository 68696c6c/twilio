package twilio_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type TwilioClient interface {
	SendSMS(to, message string) error
	GetRawResponse() http.Response
	GetResponse() TwilioResponse
}

type TwilioClientHTTP struct {
	accountSID   string
	authToken    string
	phoneNumber  string
	url          string
	client       *http.Client
	httpResponse *http.Response
	response     *TwilioResponse
}

func NewTwilioClient(accountSID, authToken, phoneNumber string) TwilioClient {
	return &TwilioClientHTTP{
		accountSID:  accountSID,
		authToken:   authToken,
		phoneNumber: phoneNumber,
		url:         fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID),
		client:      &http.Client{},
	}
}

type TwilioResponse struct {
	DateCreated         string
	DateUpdated         string
	Body                string
	NumSegments         string
	URI                 string
	To                  string
	APIVersion          string
	SubresourceURIs     map[string]interface{}
	SID                 string
	MessagingServiceSID *string
	Status              string
	ErrorCode           *string
	DateSent            *string
	AccountSID          string
	From                string
	NumMedia            string
	Direction           string
	Price               *string
	PriceUnit           string
	ErrorMessage        error
}

func (r TwilioResponse) String() string {
	out, err := json.Marshal(&r)
	if err != nil {
		return "failed to marshal response"
	}
	return string(out)
}

func (c *TwilioClientHTTP) SendSMS(to, message string) error {
	c.httpResponse = nil
	c.response = nil

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", c.phoneNumber)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	req, _ := http.NewRequest("POST", c.url, &msgDataReader)
	req.SetBasicAuth(c.accountSID, c.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	c.httpResponse = resp

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		response := &TwilioResponse{}
		body, err := streamToByte(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		err = json.Unmarshal(body, response)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}

		c.response = response

		return nil
	} else {
		return errors.Wrapf(err, "received error response from twilio api: %s", resp.Status)
	}
}

func (c *TwilioClientHTTP) GetRawResponse() http.Response {
	if c.httpResponse == nil {
		return http.Response{}
	}
	return *c.httpResponse
}

func (c *TwilioClientHTTP) GetResponse() TwilioResponse {
	if c.response == nil {
		return TwilioResponse{}
	}
	return *c.response
}

func streamToByte(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
