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
 * SOFTWARE.
 *
 */

package push

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/techcraftlabs/base"
)

var (
	_ base.RequestInformer = (*requestType)(nil)
)

const (
	SuccessCode = "BILLER-30-0000-S"
	FailureCode = "BILLER-30-3030-E"
)



type (
	
	requestType int
	Service interface {
		Token(ctx context.Context) (TokenResponse, error)
		Push(ctx context.Context, request Request) (PayResponse, error)
		CallbackServeHTTP(writer http.ResponseWriter, r *http.Request)
	}
	TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	PayRequest struct {
		CustomerMSISDN string `json:"CustomerMSISDN"`
		Amount         int64    `json:"Amount"`
		Remarks        string `json:"Remarks,omitempty"`
		ReferenceID    string `json:"ReferenceID"`
	}

	// payRequest This is the request expected by tigo with BillerMSISDN hooked up
	// from Config
	// PayRequest is used by Client without the need to specify BillerMSISDN
	payRequest struct {
		CustomerMSISDN string `json:"CustomerMSISDN"`
		BillerMSISDN   string `json:"BillerMSISDN"`
		Amount         float64   `json:"Amount"`
		Remarks        string `json:"Remarks,omitempty"`
		ReferenceID    string `json:"ReferenceID"`
	}

	PayResponse struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseStatus      bool   `json:"ResponseStatus"`
		ResponseDescription string `json:"ResponseDescription"`
		ReferenceID         string `json:"ReferenceID"`
		Message             string `json:"Message,omitempty"`
	}

	CallbackRequest struct {
		Status           bool   `json:"Status"`
		Description      string `json:"Description"`
		MFSTransactionID string `json:"MFSTransactionID,omitempty"`
		ReferenceID      string `json:"ReferenceID"`
		Amount           string `json:"Amount"`
	}

	CallbackResponse struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseDescription string `json:"ResponseDescription"`
		ResponseStatus      bool   `json:"ResponseStatus"`
		ReferenceID         string `json:"ReferenceID"`
	}

	Request struct {
		MSISDN      string  `json:"msisdn"`
		Amount      float64 `json:"amount"`
		Remarks     string  `json:"remarks,omitempty"`
		ReferenceID string  `json:"referenceID"`
	}


	Config struct {
		Username              string
		Password              string
		PasswordGrantType     string
		BaseURL            string
		TokenEndpoint          string
		BillerMSISDN          string
		BillerCode            string
		PushPayEndpoint          string
	}

	CallbackHandler interface {
		Handle(ctx context.Context, request CallbackRequest) (CallbackResponse, error)
	}

	CallbackHandlerFunc func(context.Context, CallbackRequest) (CallbackResponse, error)

	//Client is the client for making push pay requests
	Client struct {
		*Config
		base            *base.Client
		CallbackHandler CallbackHandler
		token           *string
		tokenExpires    time.Time
		rv base.Receiver
		rp base.Replier
	}
	
)



func NewClient(config *Config, handler CallbackHandler, opts ...ClientOption) *Client {
	client := &Client{
		Config:          config,
		CallbackHandler: handler,
		token:           new(string),
		tokenExpires:    time.Now(),
		base:            base.NewClient(),
	}

	for _, opt := range opts {
		opt(client)
	}

	lg, dm := client.base.Logger,client.base.DebugMode

	client.rp = base.NewReplier(lg,dm)
	client.rv = base.NewReceiver(lg,dm)

	return client
}

func (c *Client) SetCallbackHandler(handler CallbackHandler) {
	c.CallbackHandler = handler
}

func (handler CallbackHandlerFunc) Handle(ctx context.Context, request CallbackRequest) (CallbackResponse, error) {
	return handler(ctx, request)
}

func (c *Client) Push(ctx context.Context, request Request) (response PayResponse, err error) {
	amount := math.Floor(request.Amount * 100 / 100)
	var billPayReq = payRequest{
		CustomerMSISDN: request.MSISDN,
		BillerMSISDN:   c.BillerMSISDN,
		Amount:         amount,
		Remarks:        request.Remarks,
		ReferenceID:    fmt.Sprintf("%s%s", c.Config.BillerCode, request.ReferenceID),
	}

	token, err := c.checkToken(ctx)
	if err != nil {
		return PayResponse{}, err
	}

	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("bearer %s", token),
		"Username":      c.Config.Username,
		"Password":      c.Config.Password,
	}
	var requestOpts []base.RequestOption
	moreHeaderOpt := base.WithMoreHeaders(authHeader)
	//basicAuth := base.WithBasicAuth(c.PushConfig.Username, c.PushConfig.Password)
	requestOpts = append(requestOpts, moreHeaderOpt)

	req := base.MakeInternalRequest(c.BaseURL, c.Config.PushPayEndpoint, push, billPayReq, requestOpts...)
	_, err = c.base.Do(context.TODO(), req, &response)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (c *Client) CallbackServeHTTP(w http.ResponseWriter, r *http.Request) {

	var(
		callbackRequest CallbackRequest
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	//callbackRequest := new(CallbackRequest)
	statusCode := 200

	_,err := c.rv.Receive(ctx,callback.String(),r,&callbackRequest)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	callbackResponse, err := c.CallbackHandler.Handle(ctx, callbackRequest)

	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	var responseOpts []base.ResponseOption
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	headersOpt := base.WithMoreResponseHeaders(headers)

	responseOpts = append(responseOpts, headersOpt, base.WithResponseError(err))
	response := base.NewResponse(statusCode, callbackResponse, responseOpts...)
	c.rp.Reply(w, response)

}

func (c *Client) checkToken(ctx context.Context) (string, error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return "", err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpires.IsZero() && time.Until(c.tokenExpires) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return "", err
			}
		}
		token = *c.token
	}

	return token, nil
}

func (c *Client) Token(ctx context.Context) (TokenResponse, error) {
	var form = url.Values{}
	form.Set("username", c.Username)
	form.Set("password", c.Password)
	form.Set("grant_type", c.PasswordGrantType)

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Cache-Control": "no-cache",
	}

	var requestOptions []base.RequestOption
	headersOption := base.WithRequestHeaders(headers)
	requestOptions = append(requestOptions, headersOption)

	request := base.MakeInternalRequest(c.Config.BaseURL, c.Config.TokenEndpoint, token, form, requestOptions...)

	var tokenResponse TokenResponse

	_, err := c.base.Do(context.TODO(), request, &tokenResponse)

	if err != nil {
		return TokenResponse{}, err
	}

	token := tokenResponse.AccessToken

	c.token = &token

	//This set the value to when a new token will set above will be expired
	//the minus 10 is an overhead a margin for error.
	tokenExpiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn-10) * time.Second)
	c.tokenExpires = tokenExpiresAt

	return tokenResponse, nil

}