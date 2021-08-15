package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/techcraftlabs/tigopesa"
	"github.com/techcraftlabs/tigopesa/disburse"
	config2 "github.com/techcraftlabs/tigopesa/pkg/config"
	"github.com/techcraftlabs/tigopesa/push"
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

var _ push.CallbackHandler = (*handler)(nil)

type handler int

func (h handler) Do(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {
	panic("implement me")
}


func main() {
	timeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(),timeout)
	defer cancel()

	var opts []tigopesa.ClientOption
	debugOption := tigopesa.WithDebugMode(true)
	timeOutOption := tigopesa.WithTimeout(timeout)
	loggerOption := tigopesa.WithLogger(os.Stderr)
	contextOption := tigopesa.WithContext(ctx)
	httpOption := tigopesa.WithHTTPClient(http.DefaultClient)

	opts = append(opts,debugOption,timeOutOption,loggerOption,contextOption,httpOption)

	config := &config2.Overall{
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
	client := tigopesa.NewClient(config,queryHandler(1),payHandler(1),handler(1),opts...)

	disbReq := disburse.Request{
		ReferenceID: "",
		MSISDN:      "",
		Amount:      0,
	}

	response, err := client.Disburse(context.TODO(), disbReq.ReferenceID, disbReq.MSISDN, disbReq.Amount)

	if err != nil{
		fmt.Printf("error is %v\n",err)
	}

	fmt.Printf("response: %v\n",response)

	pushPayRequest := push.PayRequest{
		CustomerMSISDN: "",
		Amount:         0,
		Remarks:        "",
		ReferenceID:    "",
	}

	token, err := client.Token(context.TODO())
	if err != nil{
		fmt.Printf("error is %v\n",err)
	}

	fmt.Printf("response: %v\n",token)

	response2, err := client.Pay(context.TODO(),pushPayRequest)
	if err != nil{
		fmt.Printf("error is %v\n",err)
	}

	fmt.Printf("response: %v\n",response2)

	router := mux.NewRouter()

	router.HandleFunc("/pay",client.HandlePayment)
	router.HandleFunc("/name",client.HandleNameQuery)

}
