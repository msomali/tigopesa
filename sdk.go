package tigosdk

import (
	"context"
	"fmt"
	"github.com/techcraftt/tigosdk/aw"
	"github.com/techcraftt/tigosdk/pkg/conf"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/wa"
	"net/http"
)

const (
	NameQuery RequestType = iota
	Payment
	Callback
	Token
	Refund
	PushPay
	Heartbeat
	Disburse
)

var (
	_ Service = (*Client)(nil)
)

type (
	RequestType int

	Service interface {
		aw.DisburseHandler
		wa.Service
		push.Service
	}
	Client struct {
		*tigo.BaseClient
		*conf.Config
		wa   *wa.Client
		aw   *aw.Client
		push *push.Client
	}
)

func NewClient(config *conf.Config, base *tigo.BaseClient,
	handler wa.NameQueryHandler, paymentHandler wa.PaymentHandler, callbackHandler push.CallbackHandler) *Client {

	pushConf, payConf, disburseConf := config.Split()

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

	return &Client{
		BaseClient: base,
		Config:     config,
		wa:         payClient,
		aw:         disburseClient,
		push:       pushClient,
	}
}

func (client *Client) Disburse(ctx context.Context, request aw.DisburseRequest) (aw.DisburseResponse, error) {
	return client.aw.Disburse(ctx, request)
}

func (client *Client) HandleNameQuery(writer http.ResponseWriter, request *http.Request) {
	client.wa.HandleNameQuery(writer, request)
}

func (client *Client) HandlePayment(writer http.ResponseWriter, request *http.Request) {
	client.wa.HandlePayment(writer, request)
}

func (client *Client) Token(ctx context.Context) (string, error) {
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
func (client *Client) HandleRequest(ctx context.Context, requestType RequestType) http.HandlerFunc {
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	return func(writer http.ResponseWriter, request *http.Request) {
		switch requestType {
		case NameQuery:
			client.wa.HandleNameQuery(writer, request)
		case Payment:
			client.wa.HandlePayment(writer, request)
		case Callback:
			client.push.Callback(writer, request)
		default:
			http.Error(writer, "unknown request type", http.StatusInternalServerError)
		}
	}
}

//SendRequest like HandleRequest is experimental for neat and short API
//the problem with this API is type checking and conversion that you have
//to deal with while using it
func (client *Client)SendRequest(ctx context.Context,requestType RequestType,
	request interface{})(response interface{},err error){

	if request == nil && requestType != Token{
		return nil, fmt.Errorf("request can not be nil")
	}

	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	switch requestType {
	case Refund:
		refundReq, ok := request.(push.RefundRequest)
		if !ok{
			err = fmt.Errorf("invalid refund request")
			return nil, err
		}
		return client.push.Refund(ctx,refundReq)

	case Disburse:
		disburseReq, ok := request.(aw.DisburseRequest)
		if !ok{
			err = fmt.Errorf("invalid disburse request")
			return nil, err
		}
		return client.aw.Disburse(ctx,disburseReq)

	case PushPay:
		payReq, ok := request.(push.PayRequest)
		if !ok{
			err = fmt.Errorf("invalid push pay request")
			return nil, err
		}
		return client.push.Pay(ctx,payReq)

	case Heartbeat:
		healthReq, ok := request.(push.HealthCheckRequest)
		if !ok{
			err = fmt.Errorf("invalid health check request")
			return nil, err
		}
		return client.push.HeartBeat(ctx,healthReq)

	case Token:
		return client.push.Token(ctx)
	}

	return nil, fmt.Errorf("unsupported request type")
}
