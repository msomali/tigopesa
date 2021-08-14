package tigosdk

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/tigopesa/disburse"
	"github.com/techcraftlabs/tigopesa/pkg/conf"
	"github.com/techcraftlabs/tigopesa/pkg/tigo"
	"github.com/techcraftlabs/tigopesa/push"
	"github.com/techcraftlabs/tigopesa/ussd"
	"net/http"
)

var (
	_ Service = (*Client)(nil)
)

type (
	Service interface {
		disburse.Handler
		ussd.Service
		push.Service
	}
	Client struct {
		*tigo.BaseClient
		*conf.Config
		ussd     *ussd.Client
		disburse *disburse.Client
		push     *push.Client
	}
)

func NewClient(config *conf.Config, base *tigo.BaseClient,
	handler ussd.NameQueryHandler, paymentHandler ussd.PaymentHandler, callbackHandler push.CallbackHandler) *Client {

	pushConf, payConf, disburseConf := config.Split()

	pushClient := &push.Client{
		Config:          pushConf,
		BaseClient:      base,
		CallbackHandler: callbackHandler,
	}
	payClient := &ussd.Client{
		BaseClient:       base,
		Config:           payConf,
		PaymentHandler:   paymentHandler,
		NameQueryHandler: handler,
	}

	disburseClient := &disburse.Client{
		Config:     disburseConf,
		BaseClient: base,
	}

	return &Client{
		BaseClient: base,
		Config:     config,
		ussd:       payClient,
		disburse:   disburseClient,
		push:       pushClient,
	}
}

func (client *Client) Do(ctx context.Context, referenceId, msisdn string, amount float64) (disburse.Response, error) {
	return client.disburse.Do(ctx, referenceId, msisdn, amount)
}

func (client *Client) HandleNameQuery(writer http.ResponseWriter, request *http.Request) {
	client.ussd.HandleNameQuery(writer, request)
}

func (client *Client) HandlePayment(writer http.ResponseWriter, request *http.Request) {
	client.ussd.HandlePayment(writer, request)
}

func (client *Client) Token(ctx context.Context) (push.TokenResponse, error) {
	return client.push.Token(ctx)
}

func (client *Client) Pay(ctx context.Context, request push.PayRequest) (push.PayResponse, error) {
	return client.push.Pay(ctx, request)
}

func (client *Client) Callback(writer http.ResponseWriter, r *http.Request) {
	client.push.Callback(writer, r)
}

func (client *Client) Refund(ctx context.Context, request push.RefundRequest) (push.RefundResponse, error) {
	return client.push.Refund(ctx, request)
}

func (client *Client) HeartBeat(ctx context.Context, request push.HealthCheckRequest) (push.HealthCheckResponse, error) {
	return client.push.HeartBeat(ctx, request)
}

// HandleRequest is experimental no guarantees
func (client *Client) HandleRequest(ctx context.Context, requestName tigo.RequestName) http.HandlerFunc {
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	return func(writer http.ResponseWriter, request *http.Request) {
		switch requestName {
		case tigo.NameQueryRequest:
			client.ussd.HandleNameQuery(writer, request)
		case tigo.PaymentRequest:
			client.ussd.HandlePayment(writer, request)
		case tigo.CallbackRequest:
			client.push.Callback(writer, request)
		default:
			http.Error(writer, "unknown request type", http.StatusInternalServerError)
		}
	}
}

//SendRequest like HandleRequest is experimental for neat and short API
//the problem with this API is type checking and conversion that you have
//to deal with while using it
func (client *Client) SendRequest(ctx context.Context, requestName tigo.RequestName,
	request interface{}) (response interface{}, err error) {

	if request == nil && requestName != tigo.GetTokenRequest {
		return nil, fmt.Errorf("request can not be nil")
	}

	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	switch requestName {
	case tigo.RefundRequest:
		refundReq, ok := request.(push.RefundRequest)
		if !ok {
			err = fmt.Errorf("invalid refund request")
			return nil, err
		}
		return client.push.Refund(ctx, refundReq)

	case tigo.DisburseRequest:
		disburseReq, ok := request.(disburse.Request)
		if !ok {
			err = fmt.Errorf("invalid disburse request")
			return nil, err
		}
		return client.disburse.Do(ctx, disburseReq.ReferenceID, disburseReq.MSISDN, disburseReq.Amount)

	case tigo.PushPayRequest:
		payReq, ok := request.(push.PayRequest)
		if !ok {
			err = fmt.Errorf("invalid push pay request")
			return nil, err
		}
		return client.push.Pay(ctx, payReq)

	case tigo.HealthCheckRequest:
		healthReq, ok := request.(push.HealthCheckRequest)
		if !ok {
			err = fmt.Errorf("invalid health check request")
			return nil, err
		}
		return client.push.HeartBeat(ctx, healthReq)

	case tigo.GetTokenRequest:
		return client.push.Token(ctx)
	}

	return nil, fmt.Errorf("unsupported request type")
}
