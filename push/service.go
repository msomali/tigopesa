package push

import (
	"context"
)

type PushPay interface {
	// GetToken retrieve AuthToken for authenticating API requests.
	GetToken(ctx context.Context) (string, error)

	// BillPay initiate PushPay payment flow to deduct a specific amount from customer's Tigo pesa wallet.
	BillPay(ctx context.Context, request BillPayRequest) (BillPayResponse, error)

	// BillPayCallback handle all PushPay payment(s) status after customer purchase.
	BillPayCallback(ctx context.Context, request BillPayCallbackRequest) (BillPayResponse, error)

	// RefundPayment initiate payment refund and will be processed only if the payment was successful.
	RefundPayment(ctx context.Context, request RefundPaymentRequest) (RefundPaymentResponse, error)

	// HealthCheck check if Tigo Pesa PushPay API is up and running.
	HealthCheck(ctx context.Context, request HealthCheckResponse) (HealthCheckResponse, error)
}
