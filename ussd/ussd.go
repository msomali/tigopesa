package ussd

import (
	"context"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd/aw"
	"github.com/techcraftt/tigosdk/ussd/wa"
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
		NameQueryHandler wa.NameQueryHandler
		PaymentHandler   wa.PaymentHandler
		CallbackHandler  push.CallbackHandler
		waClient         *wa.Client
		awClient         *aw.Client
		pushClient       *push.PClient
	}
)

func (b *BigClient) Disburse(ctx context.Context, request aw.DisburseRequest) (aw.DisburseResponse, error) {
	return b.awClient.Disburse(ctx,request)
}

func (b *BigClient) HandleNameQuery(writer http.ResponseWriter, request *http.Request) {
	b.waClient.HandleNameQuery(writer,request)
}

func (b *BigClient) HandlePayment(writer http.ResponseWriter, request *http.Request) {
	b.waClient.HandlePayment(writer,request)
}

func (b *BigClient) Token(ctx context.Context) (string, error) {
	return b.pushClient.Token(ctx)
}

func (b *BigClient) Pay(ctx context.Context, request push.PayRequest) (push.PayResponse, error) {
	return b.pushClient.Pay(ctx,request)
}

func (b *BigClient) Callback(writer http.ResponseWriter, r *http.Request) {
	b.pushClient.Callback(writer,r)
}

func (b *BigClient) Refund(ctx context.Context, request push.RefundRequest) (push.RefundResponse, error) {
	return b.Refund(ctx,request)
}

func (b *BigClient) HeartBeat(ctx context.Context, request push.HealthCheckRequest) (push.HealthCheckResponse, error) {
	return b.HeartBeat(ctx,request)
}

func deriveConfigs(config *tigo.Config) (pushConf *push.Config, pay *wa.Config, disburse *aw.Config) {
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

	pushConf, payConf, disburseConf := deriveConfigs(config)

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
		NameQueryHandler: handler,
		PaymentHandler:   paymentHandler,
		CallbackHandler:  callbackHandler,
		waClient:         payClient,
		awClient:         disburseClient,
		pushClient:       pushClient,
	}
}
