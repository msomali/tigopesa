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

package main

import (
	"context"
	"github.com/gorilla/mux"
	tsdk "github.com/techcraftt/tigosdk"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd"
	"net/http"
	"os"
	"time"
)

type (
	// Application The details of application that utilizes tigopesasdk
	Application struct {
		*tsdk.Client
		svc *http.Server
	}
)

func main() {
	conf := tigo.Config{
		Username:                     "",
		Password:                     "",
		PasswordGrantType:            "",
		AccountName:                  "",
		AccountMSISDN:                "",
		BrandID:                      "",
		BillerCode:                   "",
		BillerMSISDN:                 "",
		ApiBaseURL:                   "",
		GetTokenRequestURL:           "",
		PushPayBillRequestURL:        "",
		PushPayReverseTransactionURL: "",
		PushPayHealthCheckURL:        "",
		AccountToWalletRequestURL:    "",
		AccountToWalletRequestPIN:    "",
		WalletToAccountRequestURL:    "",
		NameCheckRequestURL:          "",
	}

	var opts []tigo.ClientOption

	opts = append(
		opts,
		tigo.WithTimeout(time.Minute),
		tigo.WithContext(context.Background()),
		tigo.WithLogger(os.Stdout),
	)
	bc := tigo.NewBaseClient(conf, opts...)

	var names ussd.QuerySubscriberFunc
	{
		names = func(ctx context.Context, request ussd.SubscriberNameRequest) (ussd.SubscriberNameResponse, error) {
			return ussd.SubscriberNameResponse{}, nil
		}
	}

	var collector ussd.WalletToAccountFunc
	{
		collector = func(ctx context.Context, request ussd.WalletToAccountRequest) (ussd.WalletToAccountResponse, error) {
			return ussd.WalletToAccountResponse{}, nil
		}
	}
	var provider push.CallbackProvider
	{
		provider = func(ctx context.Context, request push.BillPayCallbackRequest) *push.BillPayResponse {
			return nil

		}
	}
	client := tsdk.NewClient(bc, names, collector, provider)

	router := mux.NewRouter()

	router.HandleFunc(client.NameCheckRequestURL, client.SubscriberNameHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(client.WalletToAccountRequestURL, client.WalletToAccountHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(client.AccountToWalletRequestURL, disburseHandleFunc).Methods(http.MethodPost)

	svc := &http.Server{
		Addr:              ":8090",
		Handler:           router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	app := Application{
		Client: client,
		svc:    svc,
	}

	err := app.svc.ListenAndServe()
	if err != nil {
		return
	}
}

func disburseHandleFunc(writer http.ResponseWriter, request *http.Request) {

}
