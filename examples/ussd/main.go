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
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	var opts []ussd.ClientOption
	debugOption := ussd.WithDebugMode(true)
	timeOutOption := ussd.WithTimeout(timeout)
	loggerOption := ussd.WithLogger(os.Stderr)
	contextOption := ussd.WithContext(ctx)
	httpOption := ussd.WithHTTPClient(http.DefaultClient)
	opts = append(opts, debugOption, timeOutOption, loggerOption, contextOption, httpOption)

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
