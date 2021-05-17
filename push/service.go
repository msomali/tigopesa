package push

import (
	"context"
	"encoding/json"
	"github.com/techcraftt/tigosdk/sdk"
	"log"
	"net/http"
)

// CallbackResponseProvider check and reports the status of the transaction.
// if transaction status
type CallbackResponseProvider func(context.Context, BillPayCallbackRequest) *BillPayResponse

type Service interface {
	// BillPay initiate Service payment flow to deduct a specific amount from customer's Tigo pesa wallet.
	BillPay(context.Context, BillPayRequest) (*BillPayResponse, error)

	// BillPayCallback handle push callback request.
	BillPayCallback(context.Context, *http.Request, http.ResponseWriter, CallbackResponseProvider)

	// RefundPayment initiate payment refund and will be processed only if the payment was successful.
	RefundPayment(context.Context, RefundPaymentRequest) (*RefundPaymentResponse, error)

	// HealthCheck check if Tigo Pesa Service API is up and running.
	HealthCheck(context.Context, HealthCheckRequest) (*HealthCheckResponse, error)
}

type client struct {
	*sdk.Client
}

func NewClient(c *sdk.Client) Service {
	return &client{c}
}

func NewClientFromConfig(config sdk.Config) Service {
	c, err := sdk.NewClient(config)
	if err != nil {
		log.Fatalln("failed to get authorization token error: ", err.Error())
	}

	return &client{c}
}

func (c *client) BillPay(ctx context.Context, billPaymentReq BillPayRequest) (*BillPayResponse, error) {
	var billPayResp = &BillPayResponse{}

	billPaymentReq.BillerMSISDN = c.BillerMSISDN
	req, err := c.NewRequest(http.MethodPost, c.PushPayBillRequestURL, sdk.JSONPayload, &billPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := c.SendWithAuth(ctx, req, billPayResp); err != nil {
		return nil, err
	}

	return billPayResp, nil
}

func (c *client) BillPayCallback(ctx context.Context, r *http.Request, w http.ResponseWriter, provider CallbackResponseProvider) {
	var callbackRequest BillPayCallbackRequest

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&callbackRequest)
	c.Log("Callback Request", sdk.JSONPayload, &callbackRequest)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	callbackResponse := provider(ctx, callbackRequest)
	c.Log("Callback Response", sdk.JSONPayload, &callbackResponse)

	if err := json.NewEncoder(w).Encode(callbackResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (c *client) RefundPayment(ctx context.Context, refundPaymentReq RefundPaymentRequest) (*RefundPaymentResponse, error) {
	var refundPaymentResp = &RefundPaymentResponse{}

	req, err := c.NewRequest(http.MethodPost, c.PushPayReverseTransactionURL, sdk.JSONPayload, refundPaymentReq)
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

	req, err := c.NewRequest(http.MethodPost, c.PushPayHealthCheckURL, sdk.JSONPayload, healthCheckReq)
	if err != nil {
		return nil, err
	}

	if err := c.Send(ctx, req, healthCheckResp); err != nil {
		return nil, err
	}

	return healthCheckResp, nil
}
