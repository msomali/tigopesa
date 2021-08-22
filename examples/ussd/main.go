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

package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/techcraftlabs/tigopesa/internal/term"
	"github.com/techcraftlabs/tigopesa/ussd"
	"net/http"
	"time"
)

var (
	_ ussd.PaymentHandler   = (*payHandler)(nil)
	_ ussd.NameQueryHandler = (*queryHandler)(nil)
)

type (
	payHandler   int
	queryHandler int
)

func PaymentHandler() ussd.PaymentHandleFunc {
	return func(ctx context.Context, request ussd.PayRequest) (ussd.PayResponse, error) {
		return ussd.PayResponse{}, nil
	}
}

func (q queryHandler) NameQuery(ctx context.Context, request ussd.NameRequest) (ussd.NameResponse, error) {
	panic("implement me")
}

func (p payHandler) PaymentRequest(ctx context.Context, request ussd.PayRequest) (ussd.PayResponse, error) {
	panic("implement me")
}

func main() {

	timeout := 60 * time.Second
	_, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	var opts []ussd.ClientOption
	debugOption := ussd.WithDebugMode(true)
	loggerOption := ussd.WithLogger(term.Stderr)
	httpOption := ussd.WithHTTPClient(http.DefaultClient)
	opts = append(opts, debugOption, loggerOption, httpOption)

	config := &ussd.Config{
		AccountName:   "",
		AccountMSISDN: "",
		BillerNumber:  "",
		RequestURL:    "",
		NamecheckURL:  "",
	}

	client := ussd.NewClient(config, payHandler(1), queryHandler(2), opts...)
	router := mux.NewRouter()
	router.HandleFunc("/pay", client.HandlePayment)
	router.HandleFunc("/name", client.HandleNameQuery)

}
