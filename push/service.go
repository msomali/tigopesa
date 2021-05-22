package push

import (
	"context"
	"encoding/json"
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
	BillPayCallback(context.Context)http.HandlerFunc

	// RefundPayment initiate payment refund and will be processed only if the payment was successful.
	RefundPayment(context.Context, RefundPaymentRequest) (*RefundPaymentResponse, error)

	// HealthCheck check if Tigo Pesa Service API is up and running.
	HealthCheck(context.Context, HealthCheckRequest) (*HealthCheckResponse, error)
}

type client struct {
	*Client
}

func newClient(c *Client) Service {
	return &client{c}
}

func NewClientFromConfig(config Config) Service {
	c, err := NewClient(config)
	if err != nil {
		log.Fatalln("failed to get authorization token error: ", err.Error())
	}

	return &client{c}
}

func (c *client) BillPay(ctx context.Context, billPaymentReq BillPayRequest) (*BillPayResponse, error) {
	var billPayResp = &BillPayResponse{}

	billPaymentReq.BillerMSISDN = c.BillerMSISDN
	req, err := c.NewRequest(http.MethodPost, c.PushPayBillRequestURL, JSONPayload, &billPaymentReq)
	if err != nil {
		return nil, err
	}

	if err := c.SendWithAuth(ctx, req, billPayResp); err != nil {
		return nil, err
	}

	return billPayResp, nil
}

func (c *client) BillPayCallback(ctx context.Context)http.HandlerFunc {
	var callbackRequest BillPayCallbackRequest

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&callbackRequest)
		c.Log("Callback Request", JSONPayload, &callbackRequest)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		callbackResponse := c.callbackHandler(ctx, callbackRequest)
		c.Log("Callback Response", JSONPayload, &callbackResponse)

		if err := json.NewEncoder(w).Encode(callbackResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

func (c *client) RefundPayment(ctx context.Context, refundPaymentReq RefundPaymentRequest) (*RefundPaymentResponse, error) {
	var refundPaymentResp = &RefundPaymentResponse{}

	req, err := c.NewRequest(http.MethodPost, c.PushPayReverseTransactionURL, JSONPayload, refundPaymentReq)
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

	req, err := c.NewRequest(http.MethodPost, c.PushPayHealthCheckURL, JSONPayload, healthCheckReq)
	if err != nil {
		return nil, err
	}

	if err := c.Send(ctx, req, healthCheckResp); err != nil {
		return nil, err
	}

	return healthCheckResp, nil
}
