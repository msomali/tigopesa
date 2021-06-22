package tigosdk

import (
	"context"
	"github.com/techcraftt/tigosdk/aw"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/wa"
	"net/http"
)

var (
	_ BigService = (*BigClient)(nil)
)

type (

	BigService interface {
		aw.DisburseHandler
		wa.Service
		push.PService
	}
	BigClient struct {
		*tigo.BaseClient
		*tigo.Config
		//NameQueryHandler wa.NameQueryHandler
		//PaymentHandler   wa.PaymentHandler
		//CallbackHandler  push.CallbackHandler
		wa   *wa.Client
		aw   *aw.Client
		push *push.PClient
	}
)

func (client *BigClient) Disburse(ctx context.Context, request aw.DisburseRequest) (aw.DisburseResponse, error) {
	return client.aw.Disburse(ctx,request)
}

func (client *BigClient) HandleNameQuery(writer http.ResponseWriter, request *http.Request) {
	client.wa.HandleNameQuery(writer,request)
}

func (client *BigClient) HandlePayment(writer http.ResponseWriter, request *http.Request) {
	client.wa.HandlePayment(writer,request)
}

func (client *BigClient) Token(ctx context.Context) (string, error) {
	return client.push.Token(ctx)
}

func (client *BigClient) Pay(ctx context.Context, request push.PayRequest) (push.PayResponse, error) {
	return client.push.Pay(ctx,request)
}

func (client *BigClient) Callback(writer http.ResponseWriter, r *http.Request) {
	client.push.Callback(writer,r)
}

func (client *BigClient) Refund(ctx context.Context, request push.RefundRequest) (push.RefundResponse, error) {
	return client.Refund(ctx,request)
}

func (client *BigClient) HeartBeat(ctx context.Context, request push.HealthCheckRequest) (push.HealthCheckResponse, error) {
	return client.HeartBeat(ctx,request)
}

func SplitConf(config *tigo.Config) (pushConf *push.Config, pay *wa.Config, disburse *aw.Config) {
	pushConf = &push.Config{
		Username:              config.PushUsername,
		Password:              config.PushPassword,
		PasswordGrantType:     config.PushPasswordGrantType,
		ApiBaseURL:            config.PushApiBaseURL,
		GetTokenURL:           config.PushGetTokenURL,
		BillerMSISDN:          config.PushBillerMSISDN,
		BillerCode:            config.PushBillerCode,
		PushPayURL:            config.PushPushPayURL,
		ReverseTransactionURL: config.PushReverseTransactionURL,
		HealthCheckURL:        config.PushHealthCheckURL,
	}

	pay = &wa.Config{
		AccountName:   config.PayAccountName,
		AccountMSISDN: config.PayAccountMSISDN,
		BillerNumber:  config.PayBillerNumber,
		RequestURL:    config.PayRequestURL,
		NamecheckURL:  config.PayNamecheckURL,
	}

	disburse = &aw.Config{
		AccountName:   config.DisburseAccountName,
		AccountMSISDN: config.DisburseAccountMSISDN,
		BrandID:       config.DisburseBrandID,
		PIN:           config.DisbursePIN,
		RequestURL:    config.DisburseRequestURL,
	}

	return
}


func NewPClient(config *tigo.Config, base *tigo.BaseClient,
	handler wa.NameQueryHandler, paymentHandler wa.PaymentHandler, callbackHandler push.CallbackHandler) *BigClient {

	pushConf, payConf, disburseConf := SplitConf(config)

	pushClient := &push.PClient{
		Config:          pushConf,
		BaseClient:      base,
		CallbackHandler: callbackHandler,
	}
	payClient := &wa.Client{
		BaseClient:       base,
		Config:           payConf,
		PaymentHandler:   paymentHandler,
		NameQueryHandler: handler,
	}

	disburseClient := &aw.Client{
		Config:     disburseConf,
		BaseClient: base,
	}

	return &BigClient{
		BaseClient:       base,
		Config:           config,
		wa:               payClient,
		aw:               disburseClient,
		push:             pushClient,
	}
}
