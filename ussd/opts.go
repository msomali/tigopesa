package ussd

import (
	"context"
	"io"
	"net/http"
	"time"
)

// ClientOption is a setter func to set Client details like
// Timeout, context, Http and Logger
type ClientOption func(client *Client)

// WithContext set the context to be used by Client in its ops
// this unset the default value which is context.TODO()
// This context value is mostly used by Handlers
func WithContext(ctx context.Context) ClientOption {
	return func(client *Client) {
		client.Ctx = ctx
	}
}

// WithTimeout used to set the Timeout used by handlers like sending requests to
// Tigo Gateway and back in case of DISBURSE_REQUEST or to set the max time for
// handlers QuerySubscriberFunc and WalletToAccountFunc while handling requests from tigo
// the default value is 1 minute
func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.Timeout = timeout
	}
}

// WithDebugMode set debug mode to true or false
func WithDebugMode(debugMode bool) ClientOption {
	return func(client *Client) {
		client.DebugMode = debugMode

	}
}

// WithLogger set a Logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out io.Writer) ClientOption {
	return func(client *Client) {
		if out == nil {
			return
		}
		client.Logger = out
	}
}

// WithHTTPClient when called unset the present http.Client and replace it
// with c. In case user tries to pass a nil value referencing the pkg
// i.e WithHTTPClient(nil), it will be ignored and the pkg wont be replaced
// Note: the new pkg Transport will be modified. It will be wrapped by another
// middleware that enables pkg to
func WithHTTPClient(httpClient *http.Client) ClientOption {

	// TODO check if its really necessary to set the default Timeout to 1 minute

	return func(client *Client) {
		if httpClient == nil {
			return
		}

		client.Http = httpClient
	}
}
