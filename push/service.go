package push

import (
	"context"
	"github.com/techcraftt/tigosdk"
	"log"
	"net/http"
)

// BillPayGetter provide interface for verification of callback request.
type BillPayGetter interface {
	GetBillPay(ctx context.Context, referenceID string) BillPayRequest
}

type Service interface {
	// BillPay initiate Service payment flow to deduct a specific amount from customer's Tigo pesa wallet.
	BillPay(context.Context, BillPayRequest) (*BillPayResponse, error)

	// BillPayCallback ...
	//BillPayCallback(context.Context, *http.Request, http.ResponseWriter, BillPayGetter) BillPayResponse

	// RefundPayment initiate payment refund and will be processed only if the payment was successful.
	RefundPayment(context.Context, RefundPaymentRequest) (*RefundPaymentResponse, error)

	// HealthCheck check if Tigo Pesa Service API is up and running.
	HealthCheck(context.Context, HealthCheckRequest) (*HealthCheckResponse, error)
}

type client struct {
	*tigosdk.Client
}

func NewClient(c *tigosdk.Client) Service {
	return &client{c}
}

func NewClientFromConfig(config tigosdk.Config) Service {
	c, err := tigosdk.NewClient(config)
	if err != nil {
		log.Fatalln("failed to get authorization token error: ", err.Error())
	}

	return &client{c}
}

func (c *client) BillPay(ctx context.Context, billPaymentReq BillPayRequest) (*BillPayResponse, error) {
	var billPayResp = &BillPayResponse{}

	billPaymentReq.BillerMSISDN = c.BillerMSISDN
	req, err := c.NewRequest(http.MethodPost, c.PushPayBillRequestURL, tigosdk.JSONRequest, &billPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := c.SendWithAuth(ctx, req, billPayResp); err != nil {
		return nil, err
	}

	return billPayResp, nil
}

func (c *client) BillPayCallback(ctx context.Context, billPayResp BillPayResponse) error {
	//todo : change implementation to support http handler
	return nil
}

func (c *client) RefundPayment(ctx context.Context, refundPaymentReq RefundPaymentRequest) (*RefundPaymentResponse, error) {
	var refundPaymentResp = &RefundPaymentResponse{}

	req, err := c.NewRequest(http.MethodPost, c.PushPayReverseTransactionURL, tigosdk.JSONRequest, refundPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := c.Send(ctx, req, refundPaymentResp); err != nil {
		return nil, err
	}

	return refundPaymentResp, nil
}

func (c *client) HealthCheck(ctx context.Context, healthCheckReq HealthCheckRequest) (*HealthCheckResponse, error) {
	var healthCheckResp = &HealthCheckResponse{}

	req, err := c.NewRequest(http.MethodPost, c.PushPayHealthCheckURL, tigosdk.JSONRequest, healthCheckReq)
	if err != nil {
		return nil, err
	}

	if err := c.Send(ctx, req, healthCheckResp); err != nil {
		return nil, err
	}

	return healthCheckResp, nil
}
