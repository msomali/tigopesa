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
	"fmt"
	"github.com/techcraftlabs/tigopesa"
	"github.com/techcraftlabs/tigopesa/pkg/config"
	"github.com/techcraftlabs/tigopesa/push"
	"github.com/techcraftlabs/tigopesa/ussd"
	"net/http"
	"os"
)

func NameQueerer() ussd.NameQueryFunc {
	return func(ctx context.Context, request ussd.NameRequest) (ussd.NameResponse, error) {
		return ussd.NameResponse{},nil
	}
}

func PayHandler()ussd.PaymentHandleFunc{
	return func(ctx context.Context, request ussd.PayRequest) (ussd.PayResponse, error) {
		return ussd.PayResponse{}, nil
	}
}

func Callbacker()push.CallbackHandlerFunc{
	return func(ctx context.Context, request push.CallbackRequest) (response push.CallbackResponse, err error) {
		return
	}
}

func main() {
	bigConfig := &config.Overall{
		PayAccountName:            "",
		PayAccountMSISDN:          "",
		PayBillerNumber:           "",
		PayRequestURL:             "",
		PayNamecheckURL:           "",
		DisburseAccountName:       "",
		DisburseAccountMSISDN:     "",
		DisburseBrandID:           "",
		DisbursePIN:               "",
		DisburseRequestURL:        "",
		PushUsername:              "",
		PushPassword:              "",
		PushPasswordGrantType:     "",
		PushApiBaseURL:            "",
		PushGetTokenURL:           "",
		PushBillerMSISDN:          "",
		PushBillerCode:            "",
		PushPayURL:                "",
		PushReverseTransactionURL: "",
		PushHealthCheckURL:        "",
	}
	var opts []tigopesa.ClientOption
	loggerOpt := tigopesa.WithLogger(os.Stderr)
	debugOpt :=tigopesa.WithDebugMode(true)
	httpClient := tigopesa.WithHTTPClient(http.DefaultClient)
	opts = append(opts,loggerOpt,debugOpt,httpClient)
	client := tigopesa.NewClient(bigConfig,NameQueerer(),PayHandler(),Callbacker(),opts...)

	//client.Disburse()
	//client.Pay()
	//client.

	fmt.Printf("client : %v\n",client)
}
