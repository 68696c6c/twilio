package twilio

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

type SMSClient interface {
	Send(to, message string) (Response, error)
}

type Client struct {
	accountSID  string
	authToken   string
	phoneNumber string
	url         string
	client      *http.Client
}

func NewClient(accountSID, authToken, phoneNumber string) SMSClient {
	return &Client{
		accountSID:  accountSID,
		authToken:   authToken,
		phoneNumber: phoneNumber,
		url:         fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID),
		client:      &http.Client{},
	}
}

type Response struct {
	HTTPResponse        *http.Response `json:"-"`
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

func (r Response) String() string {
	out, err := json.Marshal(&r)
	if err != nil {
		return "failed to marshal response"
	}
	return string(out)
}

func (c *Client) Send(to, message string) (Response, error) {
	data := url.Values{}
	data.Set("To", to)
	data.Set("From", c.phoneNumber)
	data.Set("Body", message)
	dataReader := *strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", c.url, &dataReader)
	if err != nil {
		return Response{}, errors.Wrap(err, "failed to build request")
	}

	req.SetBasicAuth(c.accountSID, c.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, errors.Wrap(err, "failed to send request")
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		response := Response{}
		body, err := streamToByte(resp.Body)
		if err != nil {
			return Response{}, errors.Wrap(err, "failed to read response body")
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return Response{}, errors.Wrap(err, "failed to unmarshal response body")
		}

		response.HTTPResponse = resp

		return response, nil
	} else {
		return Response{}, errors.Wrapf(err, "received error response from twilio api: %s", resp.Status)
	}
}

func streamToByte(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
