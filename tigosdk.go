package tigosdk

import (
	"context"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd"
	"net/http"
)

var _ Service = (*client)(nil)

const (
	SYNC_LOOKUP_RESPONSE = "SYNC_LOOKUP_RESPONSE"
	SYNC_BILLPAY_RESPONSE = "SYNC_BILLPAY_RESPONSE"
)

//Configs contains details of TigoPesa integration
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
	A2WReqURL         string `json:"a2w_url"`
}

type Service interface {
	ussd.Service
	push.Service
}

type client struct {
	Conf             Configs
	HTTPClient       *http.Client
	NameCheckHandler ussd.NameCheckFunc
	CallbackHandler  push.CallbackHandlerFunc
}

func NewTigoClient(httpClient *http.Client, nameCheckHandler ussd.NameCheckFunc, callbackHandler push.CallbackHandlerFunc, conf Configs) Service {
	return &client{
		Conf:             conf,
		HTTPClient:       httpClient,
		NameCheckHandler: nameCheckHandler,
		CallbackHandler:  callbackHandler,
	}
}

func (c client) QuerySubscriberName(ctx context.Context, req ussd.SubscriberNameRequest) (resp ussd.SubscriberNameResponse, err error) {
	nameCheckReq := ussd.NameCheckRequest{
		Msisdn:              req.Msisdn,
		CompanyName:         req.CompanyName,
		CustomerReferenceID: req.CustomerReferenceID,
	}
	re, err := c.NameCheckHandler(context.TODO(), nameCheckReq)

	if err != nil {
		return resp,err
	}

	resp = ussd.SubscriberNameResponse{
		Type:      SYNC_LOOKUP_RESPONSE,
		Result:    re.Result,
		ErrorCode: re.ErrorCode,
		ErrorDesc: re.ErrorDesc,
		Msisdn:    re.Msisdn,
		Flag:      re.Flag,
		Content:   re.Content,
	}

	return resp, nil
}

func (c client) WalletToAccount(ctx context.Context, req ussd.W2ARequest) (resp ussd.W2AResponse, err error) {
	resp = ussd.W2AResponse{
		Type:             SYNC_BILLPAY_RESPONSE,
		TxnID:            req.TxnID,
		RefID:            "dummyrefno12345",
		Result:           "TS",
		ErrorCode:        "error000",
		ErrorDescription: "Transaction Successful",
		Msisdn:           req.Msisdn,
		Flag:             "Y",
		Content:          "THE BILLPAY RESPONSE",
	}

	return
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