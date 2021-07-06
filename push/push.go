package push

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/techcraftt/tigosdk/internal"
	"github.com/techcraftt/tigosdk/pkg/tigo"
)

const (
	SuccessCode = "BILLER-30-0000-S"
	FailureCode = "BILLER-30-3030-E"
)

var (
	_ Service         = (*Client)(nil)
	_ CallbackHandler = (*CallbackHandlerFunc)(nil)
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
		Do(ctx context.Context, request CallbackRequest) (CallbackResponse, error)
	}

	CallbackHandlerFunc func(context.Context, CallbackRequest) (CallbackResponse, error)

	//Client is the client for making push pay requests
	Client struct {
		*Config
		*tigo.BaseClient
		CallbackHandler CallbackHandler
		token           string
		tokenExpires    time.Time
	}

	Service interface {
		Token(ctx context.Context) (string, error)

		Pay(ctx context.Context, request PayRequest) (PayResponse, error)

		Callback(writer http.ResponseWriter, r *http.Request)

		Refund(ctx context.Context, request RefundRequest) (RefundResponse, error)

		HeartBeat(ctx context.Context, request HealthCheckRequest) (HealthCheckResponse, error)
	}
)

func (handler CallbackHandlerFunc) Do(ctx context.Context, request CallbackRequest) (CallbackResponse, error) {
	return handler(ctx, request)
}

func (client *Client) Pay(ctx context.Context, request PayRequest) (PayResponse, error) {
	var billPayResp = PayResponse{}

	var billPayReq = payRequest{
		CustomerMSISDN: request.CustomerMSISDN,
		BillerMSISDN:   client.BillerMSISDN,
		Amount:         request.Amount,
		Remarks:        request.Remarks,
		ReferenceID:    request.ReferenceID,
	}

	var tokenStr string

	if client.token == ""{
		_, err := client.Token(ctx)

		if err != nil{
			return PayResponse{}, err
		}

		tokenStr = fmt.Sprintf("bearer %s", client.token)
	}

	//Add Auth Header
	if client.token != "" {
		if !client.tokenExpires.IsZero() && time.Until(client.tokenExpires) < (60*time.Second) {
			if _, err := client.Token(client.Ctx); err != nil {
				return billPayResp, err
			}
		}

		tokenStr = fmt.Sprintf("bearer %s", client.token)
	}


	authHeader := map[string]string{
		"Authorization": tokenStr,
	}
	var requestOpts []tigo.RequestOption
	moreHeaderOpt := tigo.WithMoreHeaders(authHeader)
	basicAuth := tigo.WithAuthHeaders(client.Username, client.Password)
	ctxOpt := tigo.WithRequestContext(ctx)
	requestOpts = append(requestOpts,ctxOpt,basicAuth,moreHeaderOpt)

	req := tigo.NewRequest(http.MethodPost,
		client.ApiBaseURL+client.PushPayURL,
		internal.JsonPayload, billPayReq,
		requestOpts...,
	)

	err := client.Send(context.TODO(), req, &billPayResp)

	if err != nil {
		return billPayResp, err
	}

	//FIXME:
	client.LogPayload(internal.JsonPayload, "PUSH RESPONSE", billPayResp)

	return billPayResp, nil

}

func (client *Client) Callback(w http.ResponseWriter, r *http.Request) {
	var callbackRequest CallbackRequest
	var callbackResponse CallbackResponse
	var response *tigo.Response
	defer func(debug bool) {
		client.Log(r, nil)
		client.LogPayload(internal.JsonPayload, "callback from tigo", &callbackRequest)
		client.LogPayload(internal.JsonPayload, "callback response", &callbackResponse)
	}(client.DebugMode)

	err := tigo.ReceiveRequest(r, internal.JsonPayload, &callbackRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	callbackResponse, err = client.CallbackHandler.Do(client.Ctx, callbackRequest)

	client.LogPayload(internal.JsonPayload, "callback response", &callbackResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var responseOpts []tigo.ResponseOption
	headers := tigo.WithDefaultJsonHeader()

	responseOpts = append(responseOpts, headers)
	response = tigo.NewResponse(200, callbackResponse, internal.JsonPayload, responseOpts...)

	_ = response.Send(w)

}

func (client *Client) Refund(ctx context.Context, refundReq RefundRequest) (RefundResponse, error) {
	var refundPaymentResp = &RefundResponse{}

	var requestOptions []tigo.RequestOption
	ctxOption := tigo.WithRequestContext(ctx)
	requestOptions = append(requestOptions, ctxOption)

	request := tigo.NewRequest(http.MethodPost,
		client.ApiBaseURL+client.GetTokenURL,
		internal.JsonPayload, refundReq,
		requestOptions...,
	)

	if err := client.Send(ctx, request, refundPaymentResp); err != nil {
		return RefundResponse{}, err
	}

	return *refundPaymentResp, nil
}

func (client *Client) HeartBeat(ctx context.Context, request HealthCheckRequest) (HealthCheckResponse, error) {
	var healthCheckResp = &HealthCheckResponse{}

	var requestOptions []tigo.RequestOption
	ctxOption := tigo.WithRequestContext(ctx)
	requestOptions = append(requestOptions, ctxOption)

	req := tigo.NewRequest(http.MethodPost, client.HealthCheckURL,
		internal.JsonPayload, request, requestOptions...)

	if err := client.Send(ctx, req, healthCheckResp); err != nil {
		return HealthCheckResponse{}, err
	}

	return *healthCheckResp, nil
}

func (client *Client) Token(ctx context.Context) (string, error) {
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

	var requestOptions []tigo.RequestOption
	ctxOption := tigo.WithRequestContext(ctx)
	headersOption := tigo.WithRequestHeaders(headers)
	requestOptions = append(requestOptions, ctxOption, headersOption)

	request := tigo.NewRequest(http.MethodPost,
		client.ApiBaseURL+client.GetTokenURL,
		payloadType, form,
		requestOptions...,
	)

	var tokenResponse TokenResponse

	err := client.Send(context.TODO(), request, &tokenResponse)

	if err != nil {
		return "", err
	}

	token := tokenResponse.AccessToken

	client.token = token

	//This set the value to when a new token will set above will be expired
	//the minus 10 is an overhead a margin for error.
	tokenExpiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn-10) * time.Second)
	client.tokenExpires = tokenExpiresAt

	return token, nil

}
