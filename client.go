package tigosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/techcraftt/tigosdk/push"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var baseURl = ""

type Client struct {
	username           string
	password           string
	authToken          string
	authTokenExpiresAt time.Time
	accountName        string
	accountMSISDN      string
	brandID            string
	billerCode         string
	client             *http.Client
	apiBaseURL         string
}

func NewClient(config Config) *Client {
	client := &Client{
		username:      config.Username,
		password:      config.Password,
		accountName:   config.AccountName,
		accountMSISDN: config.AccountMSISDN,
		brandID:       config.BrandID,
		client:        http.DefaultClient,
		billerCode:    config.BillerCode,
	}

	if _, err := client.getAuthToken(); err != nil {
		return nil
	}

	return client
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

func (c *Client) NewRequest(method, url string, payload interface{}) (*http.Request, error) {
	var (
		buf    io.Reader
		ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	)

	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	return http.NewRequestWithContext(ctx, method, url, buf)
}

func (c *Client) getAuthToken() (*push.GetTokenResponse, error) {
	var form = url.Values{}
	form.Set("username", c.username)
	form.Set("password", c.password)
	form.Set("grant_type", "password")

	var getTokenResponse = &push.GetTokenResponse{}

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiBaseURL+"/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	if err := c.Send(ctx, req, getTokenResponse); err != nil {
		return nil, err
	}

	if getTokenResponse.AccessToken != "" {
		c.authToken = getTokenResponse.AccessToken
		c.authTokenExpiresAt = time.Now().Add(time.Duration(getTokenResponse.ExpiresIn) * time.Second)
	}

	return getTokenResponse, nil
}

func (c *Client) Send(ctx context.Context, req *http.Request, v interface{}) error {
	if v == nil {
		return errors.New("v interface can not be empty")
	}

	if req.Header.Get("content-Type") == "" {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("username", c.username)
		req.Header.Set("password", c.password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
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
