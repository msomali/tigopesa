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
 * SOFTWARE.
 *
 */

package tigopesa

import (
	"context"
	"github.com/techcraftlabs/tigopesa/disburse"
	"github.com/techcraftlabs/tigopesa/internal"
	"github.com/techcraftlabs/tigopesa/pkg/config"
	"github.com/techcraftlabs/tigopesa/push"
	"github.com/techcraftlabs/tigopesa/ussd"
	"net/http"
)

var (
	_ Service = (*Client)(nil)
)

type (
	Service interface {
		disburse.Service
		ussd.Service
		push.Service
	}
	Client struct {
		*internal.BaseClient
		*config.Overall
		ussd     *ussd.Client
		disburse *disburse.Client
		push     *push.Client
	}
)

func NewClient(config *config.Overall, handler ussd.NameQueryHandler, paymentHandler ussd.PaymentHandler,
	callbackHandler push.CallbackHandler, opts ...ClientOption) *Client {

	client := &Client{
		Overall:    config,
		BaseClient: internal.NewBaseClient(),
	}

	for _, opt := range opts {
		opt(client)
	}

	pushConf, payConf, disburseConf := config.Split()

	base := &internal.BaseClient{
		HTTP:      client.HTTP,
		Logger:    client.Logger,
		DebugMode: client.DebugMode,
	}
	pushBase := base
	pushClient := makePushClient(pushConf, callbackHandler, pushBase)

	ussdBase := base
	payClient := makeUSSDClient(payConf, paymentHandler, handler, ussdBase)

	disburseBase := base
	disburseClient := makeDisburseClient(disburseConf, disburseBase)

	client.push = pushClient
	client.ussd = payClient
	client.disburse = disburseClient

	return client
}

func makePushClient(conf *push.Config, handler push.CallbackHandler, client *internal.BaseClient) *push.Client {
	logger := client.Logger
	debug := client.DebugMode
	c := client.HTTP
	var opts []push.ClientOption
	loggerOpt := push.WithLogger(logger)
	debugOpt := push.WithDebugMode(debug)
	httpClient := push.WithHTTPClient(c)
	opts = append(opts, loggerOpt, debugOpt, httpClient)

	return push.NewClient(conf, handler, opts...)
}

func makeDisburseClient(conf *disburse.Config, client *internal.BaseClient) *disburse.Client {
	logger := client.Logger
	debug := client.DebugMode
	c := client.HTTP
	var opts []disburse.ClientOption
	loggerOpt := disburse.WithLogger(logger)
	debugOpt := disburse.WithDebugMode(debug)
	httpClient := disburse.WithHTTPClient(c)
	opts = append(opts, loggerOpt, debugOpt, httpClient)

	return disburse.NewClient(conf, opts...)
}

func makeUSSDClient(conf *ussd.Config, payHandler ussd.PaymentHandler, nameHandler ussd.NameQueryHandler, client *internal.BaseClient) *ussd.Client {
	logger := client.Logger
	debug := client.DebugMode
	c := client.HTTP
	var opts []ussd.ClientOption
	loggerOpt := ussd.WithLogger(logger)
	debugOpt := ussd.WithDebugMode(debug)
	httpClient := ussd.WithHTTPClient(c)
	opts = append(opts, loggerOpt, debugOpt, httpClient)

	return ussd.NewClient(conf, payHandler, nameHandler, opts...)
}

func (client *Client) Disburse(ctx context.Context, request disburse.Request) (disburse.Response, error) {
	return client.disburse.Disburse(ctx, request)
}

func (client *Client) HandleNameQuery(writer http.ResponseWriter, request *http.Request) {
	client.ussd.HandleNameQuery(writer, request)
}

func (client *Client) HandlePayment(writer http.ResponseWriter, request *http.Request) {
	client.ussd.HandlePayment(writer, request)
}

func (client *Client) Token(ctx context.Context) (push.TokenResponse, error) {
	return client.push.Token(ctx)
}

func (client *Client) Pay(ctx context.Context, request push.PayRequest) (push.PayResponse, error) {
	return client.push.Pay(ctx, request)
}

func (client *Client) Callback(writer http.ResponseWriter, r *http.Request) {
	client.push.Callback(writer, r)
}

func (client *Client) Refund(ctx context.Context, request push.RefundRequest) (push.RefundResponse, error) {
	return client.push.Refund(ctx, request)
}

func (client *Client) HeartBeat(ctx context.Context, request push.HealthCheckRequest) (push.HealthCheckResponse, error) {
	return client.push.HeartBeat(ctx, request)
}
