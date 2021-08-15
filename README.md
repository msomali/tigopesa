# tigosdk
tigosdk is open source fully compliant tigo pesa client written in golang


## contents
1. [usage](#usage)
2. [example](#example)
2. [projects](#projects)
3. [links](#links)
4. [contributors](#contributors)
5. [sponsors](#sponsers)

## usage
```bash

go get https://github.com/techcraftlabs/tigopesa

```
## disburse example
```go

package main

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/tigopesa/disburse"
	"net/http"
	"os"
	"time"
)

func main()  {
	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(),timeout)
	defer cancel()

	var opts []disburse.ClientOption
	debugOption := disburse.WithDebugMode(true)
	timeOutOption := disburse.WithTimeout(timeout)
	loggerOption := disburse.WithLogger(os.Stderr)
	contextOption := disburse.WithContext(ctx)
	httpOption := disburse.WithHTTPClient(http.DefaultClient)

	opts = append(opts,debugOption,timeOutOption,loggerOption,contextOption,httpOption)

	config := &disburse.Config{
		AccountName:   "",
		AccountMSISDN: "",
		BrandID:       "",
		PIN:           "",
		RequestURL:    "",
	}
	client := disburse.NewClient(config,opts...)
	request := disburse.Request{
		ReferenceID: "",
		MSISDN:      "",
		Amount:      0,
	}

	response, err := client.Disburse(context.TODO(),request.ReferenceID,request.MSISDN,request.Amount)

	if err != nil{
		fmt.Printf("error occurred: %v\n",err)
	}

	fmt.Printf("the response is: %v\n",response)
}

```
### ussd example
```go

package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/techcraftlabs/tigopesa/ussd"
	"net/http"
	"os"
	"time"
)

var (
	_ ussd.PaymentHandler = (*payHandler)(nil)
	_ ussd.NameQueryHandler = (*queryHandler)(nil)
)

type (
	payHandler int
	queryHandler int
)

func (q queryHandler) NameQuery(ctx context.Context, request ussd.NameRequest) (ussd.NameResponse, error) {
	panic("implement me")
}

func (p payHandler) PaymentRequest(ctx context.Context, request ussd.PayRequest) (ussd.PayResponse, error) {
	panic("implement me")
}

func main() {

	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(),timeout)
	defer cancel()

	var opts []ussd.ClientOption
	debugOption := ussd.WithDebugMode(true)
	timeOutOption := ussd.WithTimeout(timeout)
	loggerOption := ussd.WithLogger(os.Stderr)
	contextOption := ussd.WithContext(ctx)
	httpOption := ussd.WithHTTPClient(http.DefaultClient)
	opts = append(opts,debugOption,timeOutOption,loggerOption,contextOption,httpOption)
	
	config := &ussd.Config{
		AccountName:   "",
		AccountMSISDN: "",
		BillerNumber:  "",
		RequestURL:    "",
		NamecheckURL:  "",
	}
	
	
	client := ussd.NewClient(config,payHandler(1),queryHandler(2),opts...)

	router := mux.NewRouter()
	
	router.HandleFunc("/pay",client.HandlePayment)
	router.HandleFunc("/name",client.HandleNameQuery)
	
}

```

## projects
The List of projects using this library

1. [PAYCRAFT]() - Full-Fledged Payment as a Service

## links

## contributors

1. [Bethuel Mmbaga]()
2. [Frances Ruganyumisa]()
3. [Pius Alfred]()

## sponsors

[Techcraft Technologies Co. LTD]()
