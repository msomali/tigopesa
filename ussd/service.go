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

type NameCheckHandleFunc func(context.Context, SubscriberNameRequest) SubscriberNameResponse

type WalletToAccountFunc func(ctx context.Context, request WalletToAccountRequest) WalletToAccountResponse

type Service interface {

	QuerySubscriberName(ctx context.Context, request *http.Request) (resp SubscriberNameResponse, err error)

	WalletToAccount(ctx context.Context, request *http.Request) (resp WalletToAccountResponse, err error)

	AccountToWallet(ctx context.Context, req AccountToWalletRequest) (resp AccountToWalletResponse, err error)
}
