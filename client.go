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
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd"
	"net/http"
)

var _ Service = (*Client)(nil)

type (
	Client struct {
		*tigo.BaseClient
		ussd *ussd.Client
		push *push.Client
	}

	Service interface {
		ussd.Service
		push.Service
	}
)

func (c *Client) SubscriberNameHandler(writer http.ResponseWriter, request *http.Request) {
	c.ussd.SubscriberNameHandler(writer, request)
}

func (c *Client) WalletToAccountHandler(writer http.ResponseWriter, request *http.Request) {
	c.ussd.WalletToAccountHandler(writer, request)
}

func (c *Client) AccountToWalletHandler(ctx context.Context, req ussd.AccountToWalletRequest) (resp ussd.AccountToWalletResponse, err error) {
	return c.ussd.AccountToWalletHandler(ctx, req)
}

func (c *Client) BillPay(ctx context.Context, request push.PayRequest) (*push.PayResponse, error) {
	return c.push.BillPay(ctx, request)
}

func (c *Client) BillPayCallback(ctx context.Context) http.HandlerFunc {
	return c.push.BillPayCallback(ctx)
}

func (c *Client) RefundPayment(ctx context.Context, request push.RefundRequest) (*push.RefundResponse, error) {
	return c.push.RefundPayment(ctx, request)
}

func (c *Client) HealthCheck(ctx context.Context, request push.HealthCheckRequest) (*push.HealthCheckResponse, error) {
	return c.push.HealthCheck(ctx, request)
}

func NewClient(bc *tigo.BaseClient, namesHandler ussd.QuerySubscriberFunc,
	collectionHandler ussd.WalletToAccountFunc,
	provider push.CallbackProvider) *Client {

	//todo: set push client avoid returning error in constructor
	client := &Client{
		BaseClient: bc,
		ussd:       ussd.NewClient(bc, collectionHandler, namesHandler),
		push:       push.NewClient(bc, provider),
	}

	return client
}
