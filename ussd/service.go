package ussd

import (
	"context"
)

//type TxnStatus map[string]string
//
//var (
//	Statuses = TxnStatus{
//		"200":"Success",
//		"0":"Success",
//		"410": "unable to complete transaction, amount more than maximum limit",
//	}
//)



type Service interface {

	// QuerySubscriberName is API to handle TigoPesa system Query
	// of the Customer’s Name from Partner (Third-party) – Synchronous protocol
	// req is the input from TigoPesa
	QuerySubscriberName(ctx context.Context, req SubscriberNameRequest) (resp SubscriberNameResponse, err error)

	WalletToAccount(ctx context.Context, req W2ARequest) (resp W2AResponse, err error)

	AccountToWallet(ctx context.Context, req A2WRequest) (resp A2WResponse, err error)
}
