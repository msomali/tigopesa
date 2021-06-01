package tigo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

const (
	debugKey       = "DEBUG"
	defaultTimeout = time.Minute
)


var (
	// defaultCtx is the context used by pkg when none is set
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
	//_ Service = (*Client)(nil)

	// defaultLoggerTransport is a modified http.Transport that is responsible
	// for logging all requests and responses  that a HTTPClient owned by Client
	// sent and receives
	defaultLoggerTransport = loggingTransport{
		logger: defaultWriter,
		next:   http.DefaultTransport,
	}

	// defaultHttpClient is the pkg used by library to send Http requests, specifically
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

	loggingTransport struct {
		logger io.Writer
		next   http.RoundTripper
	}

	// ClientOption is a setter func to set BaseClient details like
	// timeout, context, httpClient and logger
	ClientOption func(client *BaseClient)

	BaseClient struct {
		Config
		httpClient          *http.Client
		ctx                 context.Context
		timeout             time.Duration
		logger              io.Writer // for logging purposes
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

// WithContext set the context to be used by Client in its ops
// this unset the default value which is context.TODO()
// This context value is mostly used by Handlers
func WithContext(ctx context.Context) ClientOption {
	return func(client *BaseClient) {
		client.ctx = ctx
	}
}

// WithTimeout used to set the timeout used by handlers like sending requests to
// Tigo Gateway and back in case of Disbursement or to set the max time for
// handlers QuerySubscriberFunc and WalletToAccountFunc while handling requests from tigo
// the default value is 1 minute
func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *BaseClient) {
		client.timeout = timeout
	}
}

// WithLogger set a logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out io.Writer) ClientOption {
	return func(client *BaseClient) {
		if out == nil {
			return
		}
		client.logger = out
	}
}

// WithHTTPClient when called unset the present http.Client and replace it
// with c. In case user tries to pass a nil value referencing the pkg
// i.e WithHTTPClient(nil), it will be ignored and the pkg wont be replaced
// Note: the new pkg Transport will be modified. It will be wrapped by another
// middleware that enables pkg to
func WithHTTPClient(c *http.Client) ClientOption {

	// TODO check if its really necessary to set the default timeout to 1 minute
	//if c.Timeout == 0 {
	//	c.Timeout = defaultTimeout
	//}
	return func(client *BaseClient) {
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

func NewBaseClient(config Config, opts...ClientOption) *BaseClient {
	client := &BaseClient{
		Config:     config,
		httpClient: defaultHttpClient,
		logger:     defaultWriter,
		timeout:    defaultTimeout,
		ctx:        defaultCtx,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
