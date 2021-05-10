package push

import (
	"context"
	"github.com/techcraftt/tigosdk"
	"net/http"
)

type Service interface {
	// BillPay initiate Service payment flow to deduct a specific amount from customer's Tigo pesa wallet.
	BillPay(context.Context, BillPayRequest) (*BillPayResponse, error)

	// RefundPayment initiate payment refund and will be processed only if the payment was successful.
	RefundPayment(context.Context, RefundPaymentRequest) (*RefundPaymentResponse, error)

	// HealthCheck check if Tigo Pesa Service API is up and running.
	HealthCheck(context.Context, HealthCheckRequest) (*HealthCheckResponse, error)
}

type Push struct {
	client *tigosdk.Client
}

func New(client *tigosdk.Client) *Push {
	return &Push{client: client}
}

func NewFromConfig(config tigosdk.Config) *Push {
	return &Push{client: tigosdk.NewClient(config)}
}

func (p *Push) BillPay(ctx context.Context, billPaymentReq BillPayRequest) (*BillPayResponse, error) {
	var billPayResp = &BillPayResponse{}

	req, err := p.client.NewRequest(http.MethodPost, p.client.PushPayBillRequestURL, tigosdk.JSONRequest, billPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := p.client.SendWithAuth(ctx, req, billPayResp); err != nil {
		return nil, err
	}
	return billPayResp, nil
}

func (p *Push) BillPayCallback(ctx context.Context, billPayCallbackReq BillPayCallbackRequest) (*BillPayResponse, error) {
	//todo : change implementation to support http handler
	return nil, nil
}

func (p *Push) RefundPayment(ctx context.Context, refundPaymentReq RefundPaymentRequest) (*RefundPaymentResponse, error) {
	var refundPaymentResp = &RefundPaymentResponse{}

	req, err := p.client.NewRequest(http.MethodPost, p.client.PushPayReverseTransactionURL, tigosdk.JSONRequest, refundPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := p.client.Send(ctx, req, refundPaymentResp); err != nil {
		return nil, err
	}

	return refundPaymentResp, nil
}

func (p *Push) HealthCheck(ctx context.Context, healthCheckReq HealthCheckRequest) (*HealthCheckResponse, error) {
	var healthCheckResp = &HealthCheckResponse{}

	req, err := p.client.NewRequest(http.MethodPost, p.client.PushPayHealthCheckURL, tigosdk.JSONRequest, healthCheckReq)
	if err != nil {
		return nil, err
	}

	if err := p.client.Send(ctx, req, healthCheckResp); err != nil {
		return nil, err
	}

	return healthCheckResp, nil
}
