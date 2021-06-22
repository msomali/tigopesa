package push

import (
	"context"
	"encoding/json"
	"github.com/techcraftt/tigosdk/internal"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	_ PService = (*PClient)(nil)
)

type (

	//PClient is the client for making push pay requests
	PClient struct {
		*Config
		*tigo.BaseClient
		CallbackHandler CallbackHandler
		token           *string
		tokenExpires    time.Time
	}

	PService interface {
		Token(ctx context.Context) (string, error)

		Pay(ctx context.Context, request PayRequest) (PayResponse, error)

		Callback(writer http.ResponseWriter, r *http.Request)

		Refund(ctx context.Context, request RefundRequest) (RefundResponse, error)

		HeartBeat(ctx context.Context, request HealthCheckRequest) (HealthCheckResponse, error)
	}
)

func (client *PClient) Pay(ctx context.Context, request PayRequest) (PayResponse, error) {
	panic("implement me")
}

func (client *PClient) Callback(ctx context.Context) http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request) {
		var callbackRequest CallbackRequest

		var callbackResponse *CallbackResponse

		defer func(debugMode bool) {
			if debugMode {
				client.Log(r, nil)
				client.LogPayload(internal.JsonPayload, "Callback Response", &callbackResponse)
			}
		}(client.DebugMode)

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

		*callbackResponse, err = client.CallbackHandler.Do(ctx, callbackRequest)

		if err := json.NewEncoder(w).Encode(callbackResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

func (client *PClient) Refund(ctx context.Context, refundReq RefundRequest) (RefundResponse, error) {
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

func (client *PClient) HeartBeat(ctx context.Context, request HealthCheckRequest) (HealthCheckResponse, error) {
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

func (client *PClient) Token(ctx context.Context) (string, error) {
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
		payloadType, strings.NewReader(form.Encode()),
		requestOptions...,
	)
	//req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.ApiBaseURL+client.GetTokenURL, strings.NewReader(form.Encode()))
	//
	//
	//if err != nil {
	//	return "", err
	//}
	//
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("Cache-Control", "no-cache")
	//
	//res, err := client.HttpClient.Do(req)
	//if err != nil {
	//	return "", err
	//}
	//
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//
	//	}
	//}(res.Body)
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return "", err
	//}

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
