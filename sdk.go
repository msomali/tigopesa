package tigopesa

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/tigopesa/disburse"
	"github.com/techcraftlabs/tigopesa/internal"
	"github.com/techcraftlabs/tigopesa/pkg/config"
	"github.com/techcraftlabs/tigopesa/push"
	"github.com/techcraftlabs/tigopesa/ussd"
	"net/http"
)

var (
	_ Service = (*Client)(nil)
)

type (
	Service interface {
		disburse.Service
		ussd.Service
		push.Service
	}
	Client struct {
		*internal.BaseClient
		*config.Overall
		ussd     *ussd.Client
		disburse *disburse.Client
		push     *push.Client
	}
)

func NewClient(config *config.Overall,handler ussd.NameQueryHandler, paymentHandler ussd.PaymentHandler,
	callbackHandler push.CallbackHandler, opts ...ClientOption) *Client {

	client := &Client{
		Overall:    config,
	}

	client.Logger = defaultWriter
	client.Ctx = defaultCtx
	client.DebugMode = false
	client.Timeout = defaultTimeout
	client.Http = defaultHttpClient

	for _, opt := range opts {
		opt(client)
	}

	pushConf, payConf, disburseConf := config.Split()

	base := &internal.BaseClient{
		Http:      client.Http,
		Ctx:       client.Ctx,
		Timeout:   client.Timeout,
		Logger:    client.Logger,
		DebugMode: client.DebugMode,
	}

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

	client.push = pushClient
	client.ussd = payClient
	client.disburse = disburseClient

	return client
}

func (client *Client) Disburse(ctx context.Context, referenceId, msisdn string, amount float64) (disburse.Response, error) {
	return client.disburse.Disburse(ctx, referenceId, msisdn, amount)
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

// handleRequest is experimental no guarantees
func (client *Client) handleRequest(ctx context.Context, requestName internal.RequestName) http.HandlerFunc {
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	return func(writer http.ResponseWriter, request *http.Request) {
		switch requestName {
		case internal.NameQueryRequest:
			client.ussd.HandleNameQuery(writer, request)
		case internal.PaymentRequest:
			client.ussd.HandlePayment(writer, request)
		case internal.CallbackRequest:
			client.push.Callback(writer, request)
		default:
			http.Error(writer, "unknown request type", http.StatusInternalServerError)
		}
	}
}

//sendRequest like handleRequest is experimental for neat and short API
//the problem with this API is type checking and conversion that you have
//to deal with while using it
func (client *Client) sendRequest(ctx context.Context, requestName internal.RequestName,
	request interface{}) (response interface{}, err error) {

	if request == nil && requestName != internal.GetTokenRequest {
		return nil, fmt.Errorf("request can not be nil")
	}

	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	switch requestName {
	case internal.RefundRequest:
		refundReq, ok := request.(push.RefundRequest)
		if !ok {
			err = fmt.Errorf("invalid refund request")
			return nil, err
		}
		return client.push.Refund(ctx, refundReq)

	case internal.DISBURSE_REQUEST:
		disburseReq, ok := request.(disburse.Request)
		if !ok {
			err = fmt.Errorf("invalid disburse request")
			return nil, err
		}
		return client.disburse.Disburse(ctx, disburseReq.ReferenceID, disburseReq.MSISDN, disburseReq.Amount)

	case internal.PushPayRequest:
		payReq, ok := request.(push.PayRequest)
		if !ok {
			err = fmt.Errorf("invalid push pay request")
			return nil, err
		}
		return client.push.Pay(ctx, payReq)

	case internal.HealthCheckRequest:
		healthReq, ok := request.(push.HealthCheckRequest)
		if !ok {
			err = fmt.Errorf("invalid health check request")
			return nil, err
		}
		return client.push.HeartBeat(ctx, healthReq)

	case internal.GetTokenRequest:
		return client.push.Token(ctx)
	}

	return nil, fmt.Errorf("unsupported request type")
}
