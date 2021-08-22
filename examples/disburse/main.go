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
