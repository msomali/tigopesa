package ussd

import (
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd/aw"
	"github.com/techcraftt/tigosdk/ussd/wa"
)

type (
	BigClient struct {
		*tigo.BaseClient
		*tigo.Config
		NameQueryHandler wa.NameQueryHandler
		PaymentHandler   wa.PaymentHandler
		CallbackHandler  push.CallbackHandler
		waClient         *wa.Client
		awClient         *aw.Client
		pushClient       *push.Client
	}
)

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

	pushClient := &push.Client{
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
