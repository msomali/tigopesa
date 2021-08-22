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
	"net/http"
	"net/url"
	"time"

	"github.com/techcraftlabs/tigopesa/internal"
)

const (
	SuccessCode = "BILLER-30-0000-S"
	FailureCode = "BILLER-30-3030-E"
)

var (

)

type (
	TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	PayRequest struct {
		CustomerMSISDN string `json:"CustomerMSISDN"`
		Amount         int    `json:"Amount"`
		Remarks        string `json:"Remarks,omitempty"`
		ReferenceID    string `json:"ReferenceID"`
	}

	// payRequest This is the request expected by tigo with BillerMSISDN hooked up
	// from Config
	// PayRequest is used by Client without the need to specify BillerMSISDN
	payRequest struct {
		CustomerMSISDN string `json:"CustomerMSISDN"`
		BillerMSISDN   string `json:"BillerMSISDN"`
		Amount         int    `json:"Amount"`
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

	RefundRequest struct {
		CustomerMSISDN      string `json:"CustomerMSISDN"`
		ChannelMSISDN       string `json:"ChannelMSISDN"`
		ChannelPIN          string `json:"ChannelPIN"`
		Amount              int    `json:"Amount"`
		MFSTransactionID    string `json:"MFSTransactionID"`
		ReferenceID         string `json:"ReferenceID"`
		PurchaseReferenceID string `json:"PurchaseReferenceID"`
	}

	RefundResponse struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseStatus      bool   `json:"ResponseStatus"`
		ResponseDescription string `json:"ResponseDescription"`
		DMReferenceID       string `json:"DMReferenceID"`
		ReferenceID         string `json:"ReferenceID"`
		MFSTransactionID    string `json:"MFSTransactionID,omitempty"`
		Message             string `json:"Message,omitempty"`
	}

	HealthCheckResponse struct {
		ReferenceID string `json:"ReferenceID"`
		Description string `json:"Description"`
	}
	HealthCheckRequest struct {
		ReferenceID string `json:"ReferenceID"`
	}

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

	CallbackHandler interface {
		Respond(ctx context.Context, request CallbackRequest) (CallbackResponse, error)
	}

	CallbackHandlerFunc func(context.Context, CallbackRequest) (CallbackResponse, error)

	//Client is the client for making push pay requests
	Client struct {
		*Config
		*internal.BaseClient
		CallbackHandler CallbackHandler
		token           string
		tokenExpires    time.Time
	}

	Service interface {
		Token(ctx context.Context) (TokenResponse, error)

		Pay(ctx context.Context, request PayRequest) (PayResponse, error)

		Callback(writer http.ResponseWriter, r *http.Request)

		Refund(ctx context.Context, request RefundRequest) (RefundResponse, error)

		HeartBeat(ctx context.Context, request HealthCheckRequest) (HealthCheckResponse, error)
	}
)

func NewClient(config *Config, handler CallbackHandler, opts ...ClientOption) *Client {
	client := &Client{
		Config:          config,
		CallbackHandler: handler,
		token:           "",
		tokenExpires:    time.Now(),
		BaseClient:      internal.NewBaseClient(),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (handler CallbackHandlerFunc) Respond(ctx context.Context, request CallbackRequest) (CallbackResponse, error) {
	return handler(ctx, request)
}

func (client *Client) Pay(ctx context.Context, request PayRequest) (response PayResponse, err error) {
	var billPayReq = payRequest{
		CustomerMSISDN: request.CustomerMSISDN,
		BillerMSISDN:   client.BillerMSISDN,
		Amount:         request.Amount,
		Remarks:        request.Remarks,
		ReferenceID:    request.ReferenceID,
	}

	var tokenStr string

	if client.token == "" {
		str, err := client.Token(ctx)
		if err != nil {
			return PayResponse{}, err
		}
		tokenStr = fmt.Sprintf("bearer %s", str.AccessToken)
	}
	//Add Auth Header
	if client.token != "" {
		if !client.tokenExpires.IsZero() && time.Until(client.tokenExpires) < (60*time.Second) {
			if _, err := client.Token(client.Ctx); err != nil {
				return response, err
			}
		}
		tokenStr = fmt.Sprintf("bearer %s", client.token)
	}

	authHeader := map[string]string{
		"Authorization": tokenStr,
	}
	var requestOpts []internal.RequestOption
	moreHeaderOpt := internal.WithMoreHeaders(authHeader)
	basicAuth := internal.WithAuthHeaders(client.Username, client.Password)
	ctxOpt := internal.WithRequestContext(ctx)
	requestOpts = append(requestOpts, ctxOpt, basicAuth, moreHeaderOpt)

	tigoRequest := internal.NewRequest(ctx, http.MethodPost,
		client.ApiBaseURL+client.PushPayURL,
		internal.JsonPayload, billPayReq,
		requestOpts...,
	)

	err = client.Send(context.TODO(), internal.PushPayRequest, tigoRequest, &response)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (client *Client) Callback(w http.ResponseWriter, r *http.Request) {

	var callbackRequest CallbackRequest
	var callbackResponse CallbackResponse
	var response *internal.Response
	var statusCode int

	statusCode = 200

	err := client.Receive(r, internal.CallbackRequest, internal.JsonPayload, &callbackRequest)

	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	callbackResponse, err = client.CallbackHandler.Respond(client.Ctx, callbackRequest)

	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	var responseOpts []internal.ResponseOption
	headers := internal.WithDefaultJsonHeader()

	responseOpts = append(responseOpts, headers, internal.WithResponseError(err))
	response = internal.NewResponse(statusCode, callbackResponse, internal.JsonPayload, responseOpts...)

	client.Reply("callback response", response, w)

}

func (client *Client) Refund(ctx context.Context, refundReq RefundRequest) (RefundResponse, error) {
	var refundPaymentResp = &RefundResponse{}

	var requestOptions []internal.RequestOption
	ctxOption := internal.WithRequestContext(ctx)
	requestOptions = append(requestOptions, ctxOption)

	request := internal.NewRequest(ctx, http.MethodPost,
		client.ApiBaseURL+client.GetTokenURL,
		internal.JsonPayload, refundReq,
		requestOptions...,
	)

	if err := client.Send(ctx, internal.RefundRequest, request, refundPaymentResp); err != nil {
		return RefundResponse{}, err
	}

	return *refundPaymentResp, nil
}

func (client *Client) HeartBeat(ctx context.Context, request HealthCheckRequest) (HealthCheckResponse, error) {
	var healthCheckResp = &HealthCheckResponse{}

	var requestOptions []internal.RequestOption
	ctxOption := internal.WithRequestContext(ctx)
	requestOptions = append(requestOptions, ctxOption)

	req := internal.NewRequest(ctx, http.MethodPost, client.HealthCheckURL,
		internal.JsonPayload, request, requestOptions...)

	if err := client.Send(ctx, internal.HealthCheckRequest, req, healthCheckResp); err != nil {
		return HealthCheckResponse{}, err
	}

	return *healthCheckResp, nil
}

func (client *Client) Token(ctx context.Context) (TokenResponse, error) {
	var form = url.Values{}
	form.Set("username", client.Username)
	form.Set("password", client.Password)
	form.Set("grant_type", client.PasswordGrantType)

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Cache-Control": "no-cache",
	}

	payloadType := internal.FormPayload

	var requestOptions []internal.RequestOption
	ctxOption := internal.WithRequestContext(ctx)
	headersOption := internal.WithRequestHeaders(headers)
	requestOptions = append(requestOptions, ctxOption, headersOption)

	request := internal.NewRequest(ctx, http.MethodPost,
		client.ApiBaseURL+client.GetTokenURL,
		payloadType, form,
		requestOptions...,
	)

	var tokenResponse TokenResponse

	err := client.Send(context.TODO(), internal.GetTokenRequest, request, &tokenResponse)

	if err != nil {
		return TokenResponse{}, err
	}

	token := tokenResponse.AccessToken

	client.token = token

	//This set the value to when a new token will set above will be expired
	//the minus 10 is an overhead a margin for error.
	tokenExpiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn-10) * time.Second)
	client.tokenExpires = tokenExpiresAt

	return tokenResponse, nil

}
