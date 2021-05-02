package tigo

import (
	"context"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd"
	"net/http"
)

var _ Service = (*client)(nil)

//Configs acontains details of TigoPesa integration
//These are configurations supplied during the integration stage
type Configs struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	PasswordGrantType string `json:"password_grant_type"`
	AccountName       string `json:"account_name"`
	AccountMSISDN     string `json:"account_msisdn"`
	BrandID           string `json:"brand_id"`
	BillerCode        string `json:"biller_code"`
	GetTokenURL       string `json:"get_token_url"`
	BillURL           string `json:"biller_payment"`
}

type Service interface {
	ussd.Service
	push.Service
}

type client struct {
	HTTPClient *http.Client
	NameCheckFunc ussd.NameCheckFunc
}

func NewTigoClient(httpClient *http.Client, nameCheckFunc ussd.NameCheckFunc) Service {
	return &client{
		HTTPClient: httpClient,
		NameCheckFunc: nameCheckFunc,
	}
}

func (c client) QuerySubscriberName(ctx context.Context, req ussd.SubscriberNameRequest) (resp ussd.SubscriberNameResponse, err error) {
	panic("implement me")
}

func (c client) WalletToAccount(ctx context.Context, req ussd.W2ARequest) (resp ussd.W2AResponse, err error) {
	panic("implement me")
}

func (c client) AccountToWallet(ctx context.Context, req ussd.A2WRequest) (resp ussd.A2WResponse, err error) {
	panic("implement me")
}

func (c client) GetToken(ctx context.Context) (string, error) {
	panic("implement me")
}

func (c client) BillPay(ctx context.Context, request push.BillPayRequest) (push.BillPayResponse, error) {
	panic("implement me")
}

func (c client) BillPayCallback(ctx context.Context, request push.BillPayCallbackRequest) (push.BillPayResponse, error) {
	panic("implement me")
}

func (c client) RefundPayment(ctx context.Context, request push.RefundPaymentRequest) (push.RefundPaymentResponse, error) {
	panic("implement me")
}

func (c client) HealthCheck(ctx context.Context, request push.HealthCheckResponse) (push.HealthCheckResponse, error) {
	panic("implement me")
}
