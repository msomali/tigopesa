package push

import (
	"context"
	"encoding/json"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (

	//PClient is the client for making push pay requests
	PClient struct {
		*Config
		*tigo.BaseClient
		CallbackHandler CallbackHandler
		token *string
		tokenExpires time.Time
	}

	PService interface {
		Token(ctx context.Context)(string,error)

		Pay(ctx context.Context, request PayRequest)(PayResponse,error)

		Callback(writer http.ResponseWriter, r *http.Request)

		Refund(ctx context.Context, request RefundRequest)(RefundResponse, error)

		HeartBeat(ctx context.Context, request HealthCheckRequest)(HealthCheckResponse,error)
	}
)

func (client *PClient) Token() (string,error) {
	var form = url.Values{}
	form.Set("username",client.Username)
	form.Set("password", client.Password)
	form.Set("grant_type", client.PasswordGrantType)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.ApiBaseURL+client.GetTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	res, err := client.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body,&tokenResponse)
	if err != nil {
		return "", err
	}

	token := tokenResponse.AccessToken

	client.token = &token

	//This set the value to when a new token will set above will be expired
	//the minus 10 is an overhead a margin for error.
	tokenExpiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn-10)*time.Second)
	client.tokenExpires = tokenExpiresAt

	return token, nil

}

