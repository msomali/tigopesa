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
	"github.com/techcraftlabs/tigopesa/push"
	"net/http"
	"os"
	"time"
)

var _ push.CallbackHandler = (*handler)(nil)

type handler int

func (h handler) Respond(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {
	panic("implement me")
}

func main() {
	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	var opts []push.ClientOption
	debugOption := push.WithDebugMode(true)
	timeOutOption := push.WithTimeout(timeout)
	loggerOption := push.WithLogger(os.Stderr)
	contextOption := push.WithContext(ctx)
	httpOption := push.WithHTTPClient(http.DefaultClient)

	opts = append(opts, debugOption, timeOutOption, loggerOption, contextOption, httpOption)

	config := &push.Config{
		Username:              "",
		Password:              "",
		PasswordGrantType:     "",
		ApiBaseURL:            "",
		GetTokenURL:           "",
		BillerMSISDN:          "",
		BillerCode:            "",
		PushPayURL:            "",
		ReverseTransactionURL: "",
		HealthCheckURL:        "",
	}

	client := push.NewClient(config, handler(1), opts...)

	pushPayRequest := push.PayRequest{
		CustomerMSISDN: "",
		Amount:         0,
		Remarks:        "",
		ReferenceID:    "",
	}

	token, err := client.Token(context.TODO())
	if err != nil {
		fmt.Printf("error is %v\n", err)
	}

	fmt.Printf("response: %v\n", token)

	response, err := client.Pay(context.TODO(), pushPayRequest)
	if err != nil {
		fmt.Printf("error is %v\n", err)
	}

	fmt.Printf("response: %v\n", response)

}
