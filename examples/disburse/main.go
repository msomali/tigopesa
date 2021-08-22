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
	"github.com/techcraftlabs/tigopesa/disburse"
	"net/http"
	"os"
	"time"
)

func main() {
	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	var opts []disburse.ClientOption
	debugOption := disburse.WithDebugMode(true)
	timeOutOption := disburse.WithTimeout(timeout)
	loggerOption := disburse.WithLogger(os.Stderr)
	contextOption := disburse.WithContext(ctx)
	httpOption := disburse.WithHTTPClient(http.DefaultClient)

	opts = append(opts, debugOption, timeOutOption, loggerOption, contextOption, httpOption)

	config := &disburse.Config{
		AccountName:   "",
		AccountMSISDN: "",
		BrandID:       "",
		PIN:           "",
		RequestURL:    "",
	}
	client := disburse.NewClient(config, opts...)
	request := disburse.Request{
		ReferenceID: "",
		MSISDN:      "",
		Amount:      0,
	}

	response, err := client.Disburse(context.TODO(), request.ReferenceID, request.MSISDN, request.Amount)

	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}

	fmt.Printf("the response is: %v\n", response)
}
