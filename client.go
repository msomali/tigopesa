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
	"net/url"
	"strings"
	"time"
)

type RequestType string

var (
	JSONRequest RequestType = "json"
	XMLRequest  RequestType = "xml"
)

//Config contains details of TigoPesa integration.
//These are configurations supplied during the integration stage.
type Config struct {
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	PasswordGrantType         string `json:"grant_type"`
	AccountName               string `json:"account_name"`
	AccountMSISDN             string `json:"account_msisdn"`
	BrandID                   string `json:"brand_id"`
	BillerCode                string `json:"biller_code"`
	GetTokenRequestURL        string `json:"token_url"`
	PushPayBillRequestURL     string `json:"bill_url"`
	AccountToWalletRequestURL string `json:"a2w_url"`
	WalletToAccountRequestURL string `json:"w2a_url"`
	NameCheckRequestURL       string `json:"name_url"`
}

type Client struct {
	Config
	authToken          string
	authTokenExpiresAt time.Time
	client             *http.Client
	apiBaseURL         string
}

func NewClient(config Config) *Client {

	client := &Client{
		Config: config,
		client: http.DefaultClient,
	}

	if _, err := client.getAuthToken(); err != nil {
		return nil
	}

	return client
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
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
			b, err := xml.Marshal(&payload)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(b)
		}
	}

	return http.NewRequestWithContext(ctx, method, url, buf)
}

func (c *Client) getAuthToken() (string, error) {
	var form = url.Values{}
	form.Set("username", c.Username)
	form.Set("password", c.Password)
	form.Set("grant_type", "password")

	var getTokenResponse = map[string]interface{}{}

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiBaseURL+"/token", strings.NewReader(form.Encode()))
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
