package tigopesa

import (
	"context"
	"github.com/techcraftlabs/base/io"
	"github.com/techcraftlabs/tigopesa/disburse"
	"github.com/techcraftlabs/tigopesa/push"
	"github.com/techcraftlabs/tigopesa/ussd"
	stdio "io"
	"net/http"
	"time"
)

var (
	_ service = (*Client)(nil)
)

type (
	service interface {
		push.Service
		disburse.Service
		ussd.Service
	}
	Client struct {
		Config    *Config
		logger    stdio.Writer
		debugMode bool
		base      *http.Client
		p         *push.Client
		u         *ussd.Client
		d         *disburse.Client
	}

	Config struct {
		Disburse *disburse.Config
		Push     *push.Config
		Ussd     *ussd.Config
	}
)

func (c *Client) Token(ctx context.Context) (push.TokenResponse, error) {
	return c.p.Token(ctx)
}

func (c *Client) Pay(ctx context.Context, request push.Request) (push.PayResponse, error) {
	return c.p.Pay(ctx, request)
}

func (c *Client) CallbackServeHTTP(writer http.ResponseWriter, r *http.Request) {
	c.p.CallbackServeHTTP(writer, r)
}

func (c *Client) Disburse(ctx context.Context, request disburse.Request) (disburse.Response, error) {
	return c.d.Disburse(ctx, request)
}

func (c *Client) NameQueryServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c.u.NameQueryServeHTTP(writer, request)
}

func (c *Client) PaymentServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c.u.PaymentServeHTTP(writer, request)
}

func NewClient(config *Config, handler push.CallbackHandler, paymentHandler ussd.PaymentHandler, queryHandler ussd.NameQueryHandler, opts ...ClientOption) *Client {
	client := new(Client)
	client = &Client{
		Config:    config,
		logger:    io.Stderr,
		debugMode: true,
		base: &http.Client{
			Timeout: time.Minute,
		},
		p: nil,
		u: nil,
		d: nil,
	}

	for _, opt := range opts {
		opt(client)
	}

	disburseConfig := config.Disburse
	pushConfig := config.Push
	ussdConfig := config.Ussd
	client.d = disburse.NewClient(disburseConfig, disburse.WithLogger(client.logger),
		disburse.WithDebugMode(client.debugMode), disburse.WithHTTPClient(client.base))
	client.u = ussd.NewClient(ussdConfig, paymentHandler, queryHandler,
		ussd.WithDebugMode(client.debugMode),
		ussd.WithLogger(client.logger),
		ussd.WithHTTPClient(client.base))
	client.p = push.NewClient(pushConfig, handler, push.WithLogger(client.logger),
		push.WithDebugMode(client.debugMode), push.WithHTTPClient(client.base))
	return client
}
