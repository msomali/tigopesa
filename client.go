package tigosdk

import (
	"context"
	"net/http"
)

type (
	RequestType int
)

const (
	NameQuery RequestType = iota
	Payment
	Callback
)

// HandleRequest is experimental no guarantees
// For reliability use SubscriberNameHandler and WalletToAccountHandler
func (client *Client) HandleRequest(ctx context.Context, requestType RequestType) http.HandlerFunc {
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	return func(writer http.ResponseWriter, request *http.Request) {
		switch requestType {
		case NameQuery:
			client.HandleNameQuery(writer, request)
		case Payment:
			client.HandlePayment(writer, request)
		case Callback:
			client.Callback(writer, request)
		default:
			http.Error(writer, "unknown request type", http.StatusInternalServerError)
		}
	}
}
