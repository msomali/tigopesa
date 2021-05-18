package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	JSONPayload PayloadType = "json"
	XMLPayload  PayloadType = "xml"
)

type (
	PayloadType string

	//Config contains details of TigoPesa integration.
	//These are configurations supplied during the integration stage.
	Config struct {
		Username                     string
		Password                     string
		PasswordGrantType            string
		AccountName                  string
		AccountMSISDN                string
		BrandID                      string
		BillerCode                   string
		BillerMSISDN                 int64
		ApiBaseURL                   string
		GetTokenRequestURL           string
		PushPayBillRequestURL        string
		PushPayReverseTransactionURL string
		PushPayHealthCheckURL        string

	}

	Client struct {
		Config
		authToken          string
		authTokenExpiresAt time.Time
		client             *http.Client
		logger             io.Writer
	}
)

// NewClient initiate new tigosdk sdk used by other services.
// Default all pretty formatted requests (in and out) and responses
// will be logged to os.Sterr to use custom logger use setLogger.
func NewClient(config Config) (*Client, error) {
	client := &Client{
		Config: config,
		client: http.DefaultClient,
		logger: os.Stderr,
	}

	if _, err := client.getAuthToken(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

// SetLogger set custom logs destination.
func (c *Client) SetLogger(out io.Writer) {
	c.logger = out
}

func (c *Client) NewRequest(method, url string, payloadType PayloadType, payload interface{}) (*http.Request, error) {
	var (
		buffer io.Reader
		ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	)

	if payload != nil {
		buf, err := marshalPayload(payloadType, payload)
		if err != nil {
			return nil, err
		}

		buffer = bytes.NewBuffer(buf)
	}

	return http.NewRequestWithContext(ctx, method, c.ApiBaseURL+url, buffer)
}

func (c *Client) getAuthToken() (string, error) {
	var form = url.Values{}
	form.Set("username", c.Username)
	form.Set("password", c.Password)
	form.Set("grant_type", "password")

	var getTokenResponse = struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int64  `json:"expires_in"`
		TokenType        string `json:"token_type"`
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
	}{}

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.ApiBaseURL+c.GetTokenRequestURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	if err := c.Send(ctx, req, &getTokenResponse); err != nil {
		return "", err
	}

	// check if response contains error.
	if getTokenResponse.Error != "" {
		return "", errors.New(fmt.Sprintf("%v: %v", getTokenResponse.Error, getTokenResponse.ErrorDescription))
	}

	if getTokenResponse.AccessToken != "" {
		c.authToken = getTokenResponse.AccessToken
		c.authTokenExpiresAt = time.Now().Add(time.Duration(getTokenResponse.ExpiresIn) * time.Second)
	}

	return c.authToken, nil
}

func (c *Client) Send(ctx context.Context, req *http.Request, v interface{}) error {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	// Restore the request body content.
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if v == nil {
		return errors.New("v interface can not be empty")
	}

	// sending json request by default
	if req.Header.Get("content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("username", c.Username)
		req.Header.Set("password", c.Password)
	}

	c.logRequest(req)
	resp, err := c.client.Do(req)
	c.logResponse(resp)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.Header.Get("Content-Type") {
	case "application/json", "application/json;charset=UTF-8":
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	case "text/xml", "text/xml;charset=UTF-8":
		if err := xml.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}

	return nil
}

func (c *Client) SendWithAuth(ctx context.Context, req *http.Request, v interface{}) error {
	if c.authToken != "" {
		if !c.authTokenExpiresAt.IsZero() && time.Until(c.authTokenExpiresAt) < (60*time.Second) {
			if _, err := c.getAuthToken(); err != nil {
				return err
			}
		}

		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.authToken))
	}

	return c.Send(ctx, req, v)
}

func (c *Client) logRequest(req *http.Request) {
	if c.logger != nil && os.Getenv("DEBUG") == "true" {
		var reqDump []byte

		if req != nil {
			reqDump, _ = httputil.DumpRequest(req, true)
		}

		c.logger.Write([]byte(fmt.Sprintf("Request: %s\n \n", string(reqDump))))
	}
}

func (c *Client) logResponse(resp *http.Response) {
	if c.logger != nil && os.Getenv("DEBUG") == "true" {
		var respDump []byte

		if resp != nil {
			respDump, _ = httputil.DumpResponse(resp, true)
		}

		c.logger.Write([]byte(fmt.Sprintf("Response: %s\n \n", string(respDump))))
	}
}

// marshalPayload returns the JSON/XML encoding of payload.
func marshalPayload(payloadType PayloadType, payload interface{}) (buf []byte, err error) {
	switch payloadType {
	case JSONPayload:
		buf, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	case XMLPayload:
		buf, err = xml.MarshalIndent(payload, "", "  ")
		if err != nil {
			return nil, err
		}
	}

	return
}

// Log write v into logger.
func (c *Client) Log(prefix string, payloadType PayloadType, payload interface{}) {
	if c.logger != nil && os.Getenv("DEBUG") == "true" {
		buf, _ := marshalPayload(payloadType, payload)

		c.logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, string(buf))))
	}
}
