/*
 * MIT License
 *
 * Copyright (c) 2021 TechCraft Technologies Co. Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE
 */

package push

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/techcraftt/tigosdk/pkg/tigo"
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

	_ Service = (*Client)(nil)
)

type (
	PayloadType string
	//
	////Config contains details of TigoPesa integration.
	////These are configurations supplied during the integration stage.
	//Config struct {
	//	Username                     string
	//	Password                     string
	//	PasswordGrantType            string
	//	AccountName                  string
	//	AccountMSISDN                string
	//	BrandID                      string
	//	BillerCode                   string
	//	BillerMSISDN                 string
	//	ApiBaseURL                   string
	//	GetTokenRequestURL           string
	//	PushPayBillRequestURL        string
	//	PushPayReverseTransactionURL string
	//	PushPayHealthCheckURL        string
	//}

	// CallbackResponder check and reports the status of the transaction.
	// if transaction status
	CallbackResponder func(context.Context, BillPayCallbackRequest) *BillPayResponse

	Service interface {
		// BillPay initiate Service payment flow to deduct a specific amount from customer's Tigo pesa wallet.
		BillPay(context.Context, BillPayRequest) (*BillPayResponse, error)

		// BillPayCallback handle push callback request.
		BillPayCallback(context.Context) http.HandlerFunc

		// RefundPayment initiate payment refund and will be processed only if the payment was successful.
		RefundPayment(context.Context, RefundPaymentRequest) (*RefundPaymentResponse, error)

		// HealthCheck check if Tigo Pesa Service API is up and running.
		HealthCheck(context.Context, HealthCheckRequest) (*HealthCheckResponse, error)
	}

	Client struct {
		tigo.Config
		authToken          string
		authTokenExpiresAt time.Time
		client             *http.Client
		logger             io.Writer
		callbackHandler    CallbackResponder
	}

	ClientOption func(client *Client)
)

func (c *Client) BillPay(ctx context.Context, billPaymentReq BillPayRequest) (*BillPayResponse, error) {
	var billPayResp = &BillPayResponse{}

	billPaymentReq.BillerMSISDN = c.BillerMSISDN
	req, err := c.NewRequest(http.MethodPost, c.PushPayBillRequestURL, JSONPayload, &billPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := c.SendWithAuth(ctx, req, billPayResp); err != nil {
		return nil, err
	}

	return billPayResp, nil
}

func (c *Client) BillPayCallback(ctx context.Context) http.HandlerFunc {
	var callbackRequest BillPayCallbackRequest

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&callbackRequest)
		c.Log("Callback Request", JSONPayload, &callbackRequest)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		callbackResponse := c.callbackHandler(ctx, callbackRequest)
		c.Log("Callback Response", JSONPayload, &callbackResponse)

		if err := json.NewEncoder(w).Encode(callbackResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (c *Client) RefundPayment(ctx context.Context, refundPaymentReq RefundPaymentRequest) (*RefundPaymentResponse, error) {
	var refundPaymentResp = &RefundPaymentResponse{}

	req, err := c.NewRequest(http.MethodPost, c.PushPayReverseTransactionURL, JSONPayload, refundPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := c.Send(ctx, req, refundPaymentResp); err != nil {
		return nil, err
	}

	return refundPaymentResp, nil
}

func (c *Client) HealthCheck(ctx context.Context, healthCheckReq HealthCheckRequest) (*HealthCheckResponse, error) {
	var healthCheckResp = &HealthCheckResponse{}

	req, err := c.NewRequest(http.MethodPost, c.PushPayHealthCheckURL, JSONPayload, healthCheckReq)
	if err != nil {
		return nil, err
	}

	if err := c.Send(ctx, req, healthCheckResp); err != nil {
		return nil, err
	}

	return healthCheckResp, nil
}

// NewClient initiate new tigosdk pkg used by other services.
// Default all pretty formatted requests (in and out) and responses
// will be logged to os.Sterr to use custom logger use setLogger.
func NewClient(config tigo.Config, provider CallbackResponder, options ...ClientOption) (*Client, error) {
	client := &Client{
		Config:          config,
		client:          http.DefaultClient,
		logger:          os.Stderr,
		callbackHandler: provider,
	}

	if _, err := client.getAuthToken(); err != nil {
		return nil, err
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

// WithLogger set a logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out io.Writer) ClientOption {
	return func(client *Client) {
		if out == nil {
			return
		}
		client.logger = out
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		if client == nil {
			return
		}

		client.client = httpClient
	}
}

//func (c *Client) SetHTTPClient(pkg *http.Client) {
//	c.pkg = pkg
//}
//
//// SetLogger set custom logs destination.
//func (c *Client) SetLogger(out io.Writer) {
//	c.logger = out
//}

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
