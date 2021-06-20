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
	"github.com/techcraftt/tigosdk/internal"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	//JSONPayload PayloadType = "json"
	//XMLPayload  PayloadType = "xml"

	_ Service = (*Client)(nil)
	_ CallbackHandler = (*CallbackHandlerFunc)(nil)
)

type (
	Config struct {
		Username              string
		Password              string
		PasswordGrantType     string
		ApiBaseURL            string
		GetTokenURL           string
		BillerMSISDN          string
		BillerCode            string
		PushPayURL            string
		ReverseTransactionURL string
		HealthCheckURL        string
	}
	//PayloadType string

	CallbackHandler interface {
		Do(ctx context.Context, request CallbackRequest)(CallbackResponse,error)
	}

	// CallbackHandleFunc check and reports the status of the transaction.
	// if transaction status
	CallbackHandleFunc func(context.Context, CallbackRequest) *CallbackResponse

	// CallbackHandlerFunc check and reports the status of the transaction.
	// if transaction status
	CallbackHandlerFunc func(context.Context, CallbackRequest) (CallbackResponse,error)

	Service interface {
		// BillPay initiate Service payment flow to deduct a specific amount from customer's Tigo pesa wallet.
		BillPay(context.Context, PayRequest) (*PayResponse, error)

		// BillPayCallback handle push callback request.
		BillPayCallback(context.Context) http.HandlerFunc

		// RefundPayment initiate payment refund and will be processed only if the payment was successful.
		RefundPayment(context.Context, RefundRequest) (*RefundResponse, error)

		// HealthCheck check if Tigo Pesa Service API is up and running.
		HealthCheck(context.Context, HealthCheckRequest) (*HealthCheckResponse, error)
	}

	Client struct {
		*Config
		*tigo.BaseClient
		authToken          string
		authTokenExpiresAt time.Time
		CallbackProvider   CallbackHandleFunc
	}
)

func (handler CallbackHandlerFunc) Do(ctx context.Context, request CallbackRequest) (CallbackResponse, error) {
	return handler(ctx, request)
}

func (c *Client) BillPay(ctx context.Context, billPaymentReq PayRequest) (*PayResponse, error) {
	var billPayResp = &PayResponse{}

	billPaymentReq.BillerMSISDN = c.BillerMSISDN
	req, err := c.NewRequest(http.MethodPost, c.PushPayURL, internal.JsonPayload, &billPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := c.SendWithAuth(ctx, req, billPayResp); err != nil {
		return nil, err
	}

	return billPayResp, nil
}

func (c *Client) BillPayCallback(ctx context.Context) http.HandlerFunc {
	var callbackRequest CallbackRequest

	return func(w http.ResponseWriter, r *http.Request) {
		var callbackResponse *CallbackResponse

		defer func(debugMode bool) {
			if debugMode{
				c.Log(r,nil)
				c.LogPayload(internal.JsonPayload,"Callback Response",&callbackResponse)
			}
		}(c.DebugMode)

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&callbackRequest)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		callbackResponse = c.CallbackProvider(ctx, callbackRequest)

		if err := json.NewEncoder(w).Encode(callbackResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (c *Client) RefundPayment(ctx context.Context, refundPaymentReq RefundRequest) (*RefundResponse, error) {
	var refundPaymentResp = &RefundResponse{}

	req, err := c.NewRequest(http.MethodPost, c.ReverseTransactionURL, internal.JsonPayload, refundPaymentReq)
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

	req, err := c.NewRequest(http.MethodPost, c.HealthCheckURL, internal.JsonPayload, healthCheckReq)
	if err != nil {
		return nil, err
	}

	if err := c.Send(ctx, req, healthCheckResp); err != nil {
		return nil, err
	}

	return healthCheckResp, nil
}

// NewClient acts as a constructor of Push Pay Client.
func NewClient(bc *tigo.BaseClient, provider CallbackHandleFunc) *Client {
	client := &Client{
		BaseClient:       bc,
		CallbackProvider: provider,
	}

	if _, err := client.getAuthToken(); err != nil {
		return nil
	}

	return client
}

func (c *Client) NewRequest(method, url string, payloadType internal.PayloadType, payload interface{}) (*http.Request, error) {
	var (
		buffer io.Reader
		ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	)

	if payload != nil {
		buf, err := internal.MarshalPayload(payloadType, payload)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.ApiBaseURL+c.GetTokenURL, strings.NewReader(form.Encode()))
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
	
	resp, err := c.HttpClient.Do(req)
	
	
	//todo log here
	go func(debugMode bool) {
		if debugMode{
			c.BaseClient.Log(req,resp)
		}
	}(c.DebugMode)

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

//func (c *Client) logRequest(req *http.Request) {
//	if c.Logger != nil && os.Getenv("DEBUG") == "true" {
//		var reqDump []byte
//
//		if req != nil {
//			reqDump, _ = httputil.DumpRequest(req, true)
//		}
//
//		c.Logger.Write([]byte(fmt.Sprintf("Request: %s\n \n", string(reqDump))))
//	}
//}
//
//func (c *Client) logResponse(resp *http.Response) {
//	if c.Logger != nil && os.Getenv("DEBUG") == "true" {
//		var respDump []byte
//
//		if resp != nil {
//			respDump, _ = httputil.DumpResponse(resp, true)
//		}
//
//		c.Logger.Write([]byte(fmt.Sprintf("Response: %s\n \n", string(respDump))))
//	}
//}

//// marshalPayload returns the JSON/XML encoding of payload.
//func marshalPayload(payloadType PayloadType, payload interface{}) (buf []byte, err error) {
//	switch payloadType {
//	case JSONPayload:
//		buf, err = json.Marshal(payload)
//		if err != nil {
//			return nil, err
//		}
//	case XMLPayload:
//		buf, err = xml.MarshalIndent(payload, "", "  ")
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return
//}
//
//// Log write v into logger.
//func (c *Client) Log(prefix string, payloadType PayloadType, payload interface{}) {
//	if c.Logger != nil && os.Getenv("DEBUG") == "true" {
//		buf, _ := marshalPayload(payloadType, payload)
//
//		c.Logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, string(buf))))
//	}
//}
