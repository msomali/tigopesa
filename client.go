package tigosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	SyncLookupResponse  = "SYNC_LOOKUP_RESPONSE"
	SyncBillPayResponse = "SYNC_BILLPAY_RESPONSE"
)

var (
	JSONRequest RequestType = "json"
	XMLRequest  RequestType = "xml"
)

type (
	RequestType string

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
		ApiBaseURL                   string
		GetTokenRequestURL           string
		PushPayBillRequestURL        string
		PushPayReverseTransactionURL string
		PushPayHealthCheckURL        string
		AccountToWalletRequestURL    string
		AccountToWalletRequestPIN    string
		WalletToAccountRequestURL    string
		NameCheckRequestURL          string
	}

	Client struct {
		Config
		authToken          string
		authTokenExpiresAt time.Time
		client             *http.Client
		logger             io.Writer
	}
)

// NewClient initiate new tigosdk client used by other services.
// Default all pretty formatted requests (in and out) and responses
// will be logged to os.Sterr to use custom logger use SetLogger.
func NewClient(config Config) *Client {
	client := &Client{
		Config: config,
		client: http.DefaultClient,
		logger: os.Stderr,
	}

	if _, err := client.getAuthToken(); err != nil {
		return nil
	}

	return client
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

// SetLogger set custom logs destination.
func (c *Client) SetLogger(out io.Writer) {
	c.logger = out
}

func (c *Client) NewRequest(method, url string, requestType RequestType, payload interface{}) (*http.Request, error) {
	var (
		buf    io.Reader
		ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	)

	if payload != nil {
		switch requestType {
		case JSONRequest:
			b, err := json.Marshal(&payload)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(b)

		case XMLRequest:
			b, err := xml.MarshalIndent(&payload, "", "  ")
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(b)
		}
	}

	return http.NewRequestWithContext(ctx, method, c.ApiBaseURL+url, buf)
}

func (c *Client) getAuthToken() (string, error) {
	var form = url.Values{}
	form.Set("username", c.Username)
	form.Set("password", c.Password)
	form.Set("grant_type", "password")

	var getTokenResponse = map[string]interface{}{}

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
	if _, ok := getTokenResponse["error"]; ok {
		return "", errors.New(fmt.Sprintf("%v: %v", getTokenResponse["error"], getTokenResponse["error_description"]))
	}

	if token, ok := getTokenResponse["access_token"]; ok {
		c.authToken = token.(string)
		expiresIn := getTokenResponse["expires_in"].(float64)
		c.authTokenExpiresAt = time.Now().Add(time.Duration(int64(expiresIn)) * time.Second)
	}

	return c.authToken, nil
}

func (c *Client) Send(ctx context.Context, req *http.Request, v interface{}) error {
	if v == nil {
		return errors.New("v interface can not be empty")
	}

	// sending json request by default
	if req.Header.Get("content-Type") == "" {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("username", c.Username)
		req.Header.Set("password", c.Password)
	} else {
		req.Header.Set("Content-Type", "text/xml")
		req.Header.Set("Connection", "keep-alive")
	}

	resp, err := c.client.Do(req)
	c.log(req, resp)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch req.Header.Get("Content-Type") {
	case "application/json":
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	case "text/xml":
		if err := xml.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) SendWithAuth(ctx context.Context, req *http.Request, v interface{}) error {
	if c.authToken != "" {
		if !c.authTokenExpiresAt.IsZero() && c.authTokenExpiresAt.Sub(time.Now()) < (60*time.Second) {
			if _, err := c.getAuthToken(); err != nil {
				return err
			}
		}

		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.authToken))
	}

	return c.Send(ctx, req, v)
}

// log will dump request and response to the logger.
func (c *Client) log(req *http.Request, resp *http.Response) {
	if c.logger != nil && os.Getenv("DEBUG") == "true" {
		var reqDump, respDump []byte

		if req != nil {
			reqDump, _ = httputil.DumpRequestOut(req, true)
		}

		if resp != nil {
			respDump, _ = httputil.DumpResponse(resp, true)
		}

		c.logger.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", string(reqDump), string(respDump))))
	}
}
