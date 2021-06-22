package push

import (
	"context"
	"fmt"
	"github.com/techcraftt/tigosdk/internal"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"net/http"
	"net/url"
	"time"
)

var (
	_ Service = (*Client)(nil)
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

	CallbackHandler interface {
		Do(ctx context.Context, request CallbackRequest) (CallbackResponse, error)
	}

	CallbackHandlerFunc func(context.Context, CallbackRequest) (CallbackResponse, error)

	//Client is the client for making push pay requests
	Client struct {
		*Config
		*tigo.BaseClient
		CallbackHandler CallbackHandler
		token           *string
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

func (client *Client) Pay(ctx context.Context, request PayRequest) (PayResponse, error) {
	var billPayResp = &PayResponse{}

	var tokenStr string

	//Add Auth Header
	if *client.token != "" {
		if !client.tokenExpires.IsZero() && time.Until(client.tokenExpires) < (60*time.Second) {
			if _, err := client.Token(client.Ctx); err != nil {
				return *billPayResp,err
			}
		}

		tokenStr = fmt.Sprintf("bearer %s",*client.token)
	}

	authHeader := map[string]string{
		"Authorization": tokenStr,
	}
	var requestOpts []tigo.RequestOption
	moreHeaderOpt := tigo.WithMoreHeaders(authHeader)
	ctxOpt := tigo.WithRequestContext(ctx)
	requestOpts = append(requestOpts, moreHeaderOpt,ctxOpt)

	req := tigo.NewRequest(http.MethodPost,
		client.PushPayURL,
		internal.JsonPayload, request,
		requestOpts...,
	)

	err := client.Send(context.TODO(), req, billPayResp)

	if err != nil {
		return *billPayResp, err
	}

	return *billPayResp, nil

}

func (client *Client) Callback(w http.ResponseWriter, r *http.Request) {
	var callbackRequest CallbackRequest

	err := tigo.ReceiveRequest(r,internal.JsonPayload, &callbackRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	callbackResponse, err := client.CallbackHandler.Do(client.Ctx, callbackRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var responseOpts []tigo.ResponseOption
	headers := tigo.WithDefaultJsonHeader()

	responseOpts = append(responseOpts, headers)
	response := tigo.NewResponse(200, callbackResponse,internal.JsonPayload, responseOpts...)

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
	//err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	token := tokenResponse.AccessToken

	client.token = &token

	//This set the value to when a new token will set above will be expired
	//the minus 10 is an overhead a margin for error.
	tokenExpiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn-10) * time.Second)
	client.tokenExpires = tokenExpiresAt

	return token, nil

}
