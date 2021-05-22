/*
 * MIT License
 *
 * Copyright (c) 2021 TechCraft Technologies Co. Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE
 */

package tigosdk

import (
	"context"
	"fmt"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

const (
	debugKey            = "DEBUG"
	defaultTimeout      = time.Minute
)

var (
	// defaultCtx is the context used by client when none is set
	// to override this one has to call WithContext method and supply
	// his her own context.Context
	defaultCtx = context.TODO()

	// defaultWriter is an io.Writer used for debugging. When debug mode is
	// set to true i.e DEBUG=true and no io.Writer is provided via
	// WithLogger method this is used.
	defaultWriter = os.Stderr

	// loggingTransport implements http.RoundTripper
	_ http.RoundTripper = (*loggingTransport)(nil)

	// Client implements Service
	_ Service = (*Client)(nil)

	// defaultLoggerTransport is a modified http.Transport that is responsible
	// for logging all requests and responses  that a HTTPClient owned by Client
	// sent and receives
	defaultLoggerTransport = loggingTransport{
		logger: defaultWriter,
		next:   http.DefaultTransport,
	}

	// defaultHttpClient is the client used by library to send Http requests, specifically
	// disbursement requests in case a user does not specify one
	defaultHttpClient = &http.Client{
		Transport: defaultLoggerTransport,
		Timeout:   defaultTimeout,
	}
)


type (
	Config struct {
		Username                     string
		Password                     string
		PasswordGrantType            string
		AccountName                  string
		AccountMSISDN                string
		BrandID                      string
		BillerCode                   string
		BillerMSISDN                 string
		ApiBaseURL                   string
		GetTokenRequestURL           string
		PushPayBillRequestURL        string
		PushPayReverseTransactionURL string
		PushPayHealthCheckURL        string
		AccountToWalletRequestURL    string
		AccountToWalletRequestPIN    string
		WalletToAccountRequestURL    string
		NameCheckRequestURL          string
	}

	Client struct {
		Config
		ussd                    ussd.Client
		push                    push.Client
		httpClient              *http.Client
		logger                  io.Writer
		timeout                 time.Duration
		ctx                     context.Context
		QuerySubscriberNameFunc ussd.QuerySubscriberFunc
		WalletToAccountFunc     ussd.WalletToAccountFunc
		CallbackResponder push.CallbackResponseProvider
	}

	ClientOption func (client *Client)

	Service interface {
		ussd.Service
		push.Service
	}


	loggingTransport struct {
		logger io.Writer
		next   http.RoundTripper
	}


)

func (l loggingTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {

	if os.Getenv(debugKey) == "true" && request != nil {
		reqDump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fprint(l.logger, fmt.Sprintf("Request %s\n", string(reqDump)))
		if err != nil {
			return nil, err
		}

	}
	defer func() {
		if response != nil && os.Getenv(debugKey) == "true" {
			respDump, err := httputil.DumpResponse(response, true)
			_, err = l.logger.Write([]byte(fmt.Sprintf("Response %s\n", string(respDump))))
			if err != nil {
				return
			}
		}
	}()
	response, err = l.next.RoundTrip(request)
	return
}

func SubscriberNameFunc(names ussd.QuerySubscriberFunc) ClientOption {
	return func(client *Client) {
		client.QuerySubscriberNameFunc = names
	}
}

func (client *Client) SubscriberNameHandler(writer http.ResponseWriter, request *http.Request) {
	client.ussd.SubscriberNameHandler(writer,request)
}

func (client *Client) WalletToAccountHandler(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (client *Client) AccountToWalletHandler(ctx context.Context, req ussd.AccountToWalletRequest) (resp ussd.AccountToWalletResponse, err error) {
	panic("implement me")
}

func (client *Client) BillPay(ctx context.Context, request push.BillPayRequest) (*push.BillPayResponse, error) {
	panic("implement me")
}

func (client *Client) BillPayCallback(ctx context.Context)http.HandlerFunc {
	panic("implement me")
}

func (client *Client) RefundPayment(ctx context.Context, request push.RefundPaymentRequest) (*push.RefundPaymentResponse, error) {
	panic("implement me")
}

func (client *Client) HealthCheck(ctx context.Context, request push.HealthCheckRequest) (*push.HealthCheckResponse, error) {
	panic("implement me")
}

func NewClient(config Config,namesHandler ussd.QuerySubscriberFunc,
	collectionHandler ussd.WalletToAccountFunc,
	callbackResponder push.CallbackResponseProvider,
	opts ...ClientOption)*Client{
	client := &Client{
		Config:                  config,
		httpClient:              defaultHttpClient,
		logger:                  defaultWriter,
		timeout:                 defaultTimeout,
		ctx:                     defaultCtx,
		QuerySubscriberNameFunc: namesHandler,
		WalletToAccountFunc:     collectionHandler,
		CallbackResponder: callbackResponder,
	}

	for _, opt := range opts {
		opt(client)
	}

	// todo: set ussd client
	client.setUSSDClient()


	// todo: set push client
	client.setPushClient()


	return client
}


// WithContext set the context to be used by Client in its ops
// this unset the default value which is context.TODO()
// This context value is mostly used by Handlers
func WithContext(ctx context.Context) ClientOption {
	return func(client *Client) {
		client.ctx = ctx
	}
}

// WithTimeout used to set the timeout used by handlers like sending requests to
// Tigo Gateway and back in case of Disbursement or to set the max time for
// handlers QuerySubscriberFunc and WalletToAccountFunc while handling requests from tigo
// the default value is 1 minute
func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.timeout = timeout
	}
}

// WithLogger set a logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out io.Writer) ClientOption {
	return func(client *Client) {
		if out == nil {
			return
		}
		client.logger = out
	}
}

// WithHTTPClient when called unset the present http.Client and replace it
// with c. In case user tries to pass a nil value referencing the client
// i.e WithHTTPClient(nil), it will be ignored and the client wont be replaced
// Note: the new client Transport will be modified. It will be wrapped by another
// middleware that enables client to
func WithHTTPClient(c *http.Client) ClientOption {

	// TODO check if its really necessary to set the default timeout to 1 minute
	//if c.Timeout == 0 {
	//	c.Timeout = defaultTimeout
	//}
	return func(client *Client) {
		if c == nil {
			return
		}

		lt := loggingTransport{
			logger: client.logger,
			next:   c.Transport,
		}

		hc := &http.Client{
			Transport:     lt,
			CheckRedirect: c.CheckRedirect,
			Jar:           c.Jar,
			Timeout:       c.Timeout,
		}
		client.httpClient = hc
	}
}

///helpers
func (client *Client) setUSSDClient()  {
	newClient := ussd.NewClient(
		ussd.Config{
			Username:                  client.Username,
			Password:                  client.Password,
			AccountName:               client.AccountName,
			AccountMSISDN:             client.AccountMSISDN,
			BrandID:                   client.AccountMSISDN,
			BillerCode:                client.BillerCode,
			AccountToWalletRequestURL: client.AccountToWalletRequestURL,
			AccountToWalletRequestPIN: client.AccountToWalletRequestPIN,
			WalletToAccountRequestURL: client.WalletToAccountRequestURL,
			NameCheckRequestURL:       client.NameCheckRequestURL,
		},
		client.WalletToAccountFunc,
		client.QuerySubscriberNameFunc,
		ussd.WithLogger(client.logger),
		ussd.WithHTTPClient(client.httpClient),
		ussd.WithContext(client.ctx),
		ussd.WithTimeout(client.timeout),
	)

	client.ussd = *newClient
}

///helpers
func (client *Client) setPushClient()  {


}



