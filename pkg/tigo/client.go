package tigo

import (
	"context"
	"fmt"
	"github.com/techcraftt/tigosdk/internal"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

const (
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

	//// loggingTransport implements http.RoundTripper
	//_ http.RoundTripper = (*loggingTransport)(nil)
	//
	//// Client implements Service
	////_ Service = (*Client)(nil)
	//
	//// defaultLoggerTransport is a modified http.Transport that is responsible
	//// for logging all requests and responses  that a HTTPClient owned by Client
	//// sent and receives
	//defaultLoggerTransport = loggingTransport{
	//	debugMode: false,
	//	logger: defaultWriter,
	//	next:   http.DefaultTransport,
	//}

	// defaultHttpClient is the pkg used by library to send Http requests, specifically
	// disbursement requests in case a user does not specify one
	defaultHttpClient = &http.Client{
		Timeout: defaultTimeout,
	}
)

type (
	Config struct {
		PayAccountName   string
		PayAccountMSISDN string
		PayBillerNumber  string
		PayRequestURL    string
		PayNamecheckURL  string


		DisburseAccountName   string
		DisburseAccountMSISDN string
		DisburseBrandID       string
		DisbursePIN           string
		DisburseRequestURL    string

		PushUsername              string
		PushPassword              string
		PushPasswordGrantType     string
		PushApiBaseURL            string
		PushGetTokenURL           string
		PushBillerMSISDN          string
		PushBillerCode            string
		PushPushPayURL            string
		PushReverseTransactionURL string
		PushHealthCheckURL        string
	}

	//loggingTransport struct {
	//	debugMode bool
	//	logger io.Writer
	//	next   http.RoundTripper
	//}

	// ClientOption is a setter func to set BaseClient details like
	// Timeout, context, HttpClient and Logger
	ClientOption func(client *BaseClient)

	BaseClient struct {
		//Config
		HttpClient *http.Client
		Ctx        context.Context
		Timeout    time.Duration
		Logger     io.Writer // for logging purposes
		DebugMode  bool
	}
)

//func (l loggingTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {
//
//	if l.debugMode && request != nil {
//		reqDump, err := httputil.DumpRequestOut(request, true)
//		if err != nil {
//			return nil, err
//		}
//		_, err = fmt.Fprint(l.logger, fmt.Sprintf("Request %s\n", string(reqDump)))
//		if err != nil {
//			return nil, err
//		}
//
//	}
//	defer func() {
//		if response != nil && l.debugMode{
//			respDump, err := httputil.DumpResponse(response, true)
//			_, err = l.logger.Write([]byte(fmt.Sprintf("Response %s\n", string(respDump))))
//			if err != nil {
//				return
//			}
//		}
//	}()
//	response, err = l.next.RoundTrip(request)
//	return
//}

// WithContext set the context to be used by Client in its ops
// this unset the default value which is context.TODO()
// This context value is mostly used by Handlers
func WithContext(ctx context.Context) ClientOption {
	return func(client *BaseClient) {
		client.Ctx = ctx
	}
}

// WithTimeout used to set the Timeout used by handlers like sending requests to
// Tigo Gateway and back in case of Disbursement or to set the max time for
// handlers QuerySubscriberFunc and WalletToAccountFunc while handling requests from tigo
// the default value is 1 minute
func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *BaseClient) {
		client.Timeout = timeout
	}
}

// WithDebugMode set debug mode to true or false
func WithDebugMode(debugMode bool) ClientOption {
	return func(client *BaseClient) {
		client.DebugMode = debugMode

	}
}

// WithLogger set a Logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out io.Writer) ClientOption {
	return func(client *BaseClient) {
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

	return func(client *BaseClient) {
		if httpClient == nil {
			return
		}

		//lt := loggingTransport{
		//	debugMode: client.DebugMode,
		//	logger: client.Logger,
		//	next:   c.Transport,
		//}
		//
		//hc := &http.Client{
		//	Transport:     lt,
		//	CheckRedirect: c.CheckRedirect,
		//	Jar:           c.Jar,
		//	Timeout:       c.Timeout,
		//}
		client.HttpClient = httpClient
	}
}

func NewBaseClient(opts ...ClientOption) *BaseClient {
	client := &BaseClient{
		//Config:     config,
		HttpClient: defaultHttpClient,
		Logger:     defaultWriter,
		Timeout:    defaultTimeout,
		Ctx:        defaultCtx,
		DebugMode:  false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (client *BaseClient)LogPayload(t internal.PayloadType,prefix string,payload interface{}){

	errors := make(chan error)
	done := make(chan bool)

	go func() {
		buf, err := internal.MarshalPayload(t, payload)
		if err != nil {
			errors <- fmt.Errorf("could not marshal the payload: %s\n", err.Error())
		}
		_, err = client.Logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, string(buf))))
		if err != nil {
			errors <- fmt.Errorf("could not print the payload: %s\n", err.Error())
		}

		done <- true
	}()

	select {
	case err := <-errors:
		fmt.Printf("error encountered: %s\n", err)
		return

	case <-done:
		fmt.Printf("done logging\n")
		return
	}
}

func (client *BaseClient)Log(request *http.Request, response *http.Response) {

	errors := make(chan error)
	done := make(chan bool)

	go func() {
		if request != nil {
			reqDump, err := httputil.DumpRequestOut(request, true)
			if err != nil {
				errors <- fmt.Errorf("could not dump request due to: %s\n", err.Error())
			}
			_, err = fmt.Fprintf(client.Logger, "request %s\n", reqDump)
			if err != nil {
				errors <- fmt.Errorf("could not print  request due to: %s\n", err.Error())
			}

			if response == nil {
				done <- true
			}
		}

		if response != nil {
			respDump, err := httputil.DumpResponse(response, true)
			if err != nil {
				errors <- fmt.Errorf("could not dump response due to: %s\n", err.Error())
			}
			_, err = fmt.Fprintf(client.Logger, "response:  %s\n", respDump)

			if err != nil {
				errors <- fmt.Errorf("could not print out response due to: %s\n", err.Error())
			}

			done <- true

		}
	}()

	select {
	case err := <-errors:
		fmt.Printf("error encountered: %s\n", err)
		return

	case <-done:
		fmt.Printf("done logging\n")
		return
	}

}
