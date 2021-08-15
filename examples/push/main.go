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
	ctx, cancel := context.WithTimeout(context.TODO(),timeout)
	defer cancel()

	var opts []push.ClientOption
	debugOption := push.WithDebugMode(true)
	timeOutOption := push.WithTimeout(timeout)
	loggerOption := push.WithLogger(os.Stderr)
	contextOption := push.WithContext(ctx)
	httpOption := push.WithHTTPClient(http.DefaultClient)

	opts = append(opts,debugOption,timeOutOption,loggerOption,contextOption,httpOption)

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
	if err != nil{
		fmt.Printf("error is %v\n",err)
	}

	fmt.Printf("response: %v\n",token)

	response, err := client.Pay(context.TODO(),pushPayRequest)
	if err != nil{
		fmt.Printf("error is %v\n",err)
	}

	fmt.Printf("response: %v\n",response)

}
