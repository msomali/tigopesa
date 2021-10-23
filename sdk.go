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
	_ Service = (*Client)(nil)
)

type (
	Service interface {
		push.Service
		disburse.Service
		ussd.Service
	}
	Client struct {
		logger stdio.Writer
		debugMode bool
		base *http.Client
		p *push.Client
		u *ussd.Client
		d *disburse.Client
	}
	
	Config struct {
		Disburse *disburse.Config
		Push *push.Config
		Ussd *ussd.Config
	}
)

func (c *Client) Token(ctx context.Context) (push.TokenResponse, error) {
	return c.p.Token(ctx)
}

func (c *Client) Push(ctx context.Context, request push.Request) (push.PayResponse, error) {
	return c.p.Push(ctx,request)
}

func (c *Client) CallbackServeHTTP(writer http.ResponseWriter, r *http.Request) {
	c.p.CallbackServeHTTP(writer,r)
}

func (c *Client) Disburse(ctx context.Context, request disburse.Request) (disburse.Response, error) {
	return c.d.Disburse(ctx,request)
}

func (c *Client) NameQueryServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c.u.NameQueryServeHTTP(writer,request)
}

func (c *Client) PaymentServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c.u.PaymentServeHTTP(writer,request)
}

func NewClient(config *Config,handler push.CallbackHandler,paymentHandler ussd.PaymentHandler, queryHandler ussd.NameQueryHandler,opts...ClientOption) *Client {
	client := new(Client)
	
	client = &Client{
		base: &http.Client{
			Timeout: time.Minute,
		},
		logger:    io.Stderr,
		debugMode: true,
	}

	for _, opt := range opts {
		opt(client)
	}
	client.d = disburse.NewClient(config.Disburse,disburse.WithLogger(client.logger),
		disburse.WithDebugMode(client.debugMode), disburse.WithHTTPClient(client.base))
	client.u = ussd.NewClient(config.Ussd,paymentHandler,queryHandler,
		ussd.WithDebugMode(client.debugMode),
		ussd.WithLogger(client.logger),
		ussd.WithHTTPClient(client.base))
	client.p = push.NewClient(config.Push,handler,push.WithLogger(client.logger),
		push.WithDebugMode(client.debugMode),push.WithHTTPClient(client.base))
	return client
}
