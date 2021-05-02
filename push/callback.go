package push

import "context"

type CallbackHandlerFunc func(context.Context, BillPayCallbackRequest)

func (f CallbackHandlerFunc) Handle(ctx context.Context, request BillPayCallbackRequest) {
	f(ctx, request)
}

type CallbackHandler interface {
	Handle(ctx context.Context, request BillPayCallbackRequest)
}
