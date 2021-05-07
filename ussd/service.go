package ussd

import (
	"context"
	"net/http"
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

type NameCheckHandleFunc func(context.Context, SubscriberNameRequest) (SubscriberNameResponse, error)

type WalletToAccountFunc func(ctx context.Context, request WalletToAccountRequest)(WalletToAccountResponse, error)

type Service interface {

	// QuerySubscriberName is API to handle TigoPesa system Query
	// of the Customer’s Name from Partner (Third-party) – Synchronous protocol
	// req is the input from TigoPesa
	QuerySubscriberName(ctx context.Context, req SubscriberNameRequest) (resp SubscriberNameResponse, err error)

	WalletToAccount(ctx context.Context, req WalletToAccountRequest) (resp WalletToAccountResponse, err error)

	QuerySubscriberNameL(ctx context.Context, request *http.Request) (resp SubscriberNameResponse, err error)

	WalletToAccountN(ctx context.Context, request *http.Request) (resp WalletToAccountResponse, err error)


	AccountToWallet(ctx context.Context, req AccountToWalletRequest) (resp AccountToWalletResponse, err error)
}
