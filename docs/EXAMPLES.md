## all

```go
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
	_ ussd.PaymentHandler   = (*payHandler)(nil)
	_ ussd.NameQueryHandler = (*queryHandler)(nil)
)

type (
	payHandler   int
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

func (h handler) Respond(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {
	panic("implement me")
}

func main() {
	timeout := 60 * time.Second
	_, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	var opts []tigopesa.ClientOption
	debugOption := tigopesa.WithDebugMode(true)
	loggerOption := tigopesa.WithLogger(os.Stderr)
	httpOption := tigopesa.WithHTTPClient(http.DefaultClient)

	opts = append(opts, debugOption, loggerOption,  httpOption)

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
	client := tigopesa.NewClient(config, queryHandler(1), payHandler(1), handler(1), opts...)

	disbReq := disburse.Request{
		ReferenceID: "",
		MSISDN:      "",
		Amount:      0,
	}

	response, err := client.Disburse(context.TODO(), disbReq.ReferenceID, disbReq.MSISDN, disbReq.Amount)

	if err != nil {
		fmt.Printf("error is %v\n", err)
	}

	fmt.Printf("response: %v\n", response)

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

	response2, err := client.Pay(context.TODO(), pushPayRequest)
	if err != nil {
		fmt.Printf("error is %v\n", err)
	}

	fmt.Printf("response: %v\n", response2)

	router := mux.NewRouter()

	router.HandleFunc("/pay", client.HandlePayment)
	router.HandleFunc("/name", client.HandleNameQuery)

}


```

## push

```go


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
	_, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	var opts []push.ClientOption
	debugOption := push.WithDebugMode(true)
	loggerOption := push.WithLogger(os.Stderr)
	httpOption := push.WithHTTPClient(http.DefaultClient)

	opts = append(opts, debugOption, loggerOption, httpOption)

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

```

## disburse

```go

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


```

## ussd

```go


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
	_, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	var opts []disburse.ClientOption
	debugOption := disburse.WithDebugMode(true)

	loggerOption := disburse.WithLogger(os.Stderr)
	httpOption := disburse.WithHTTPClient(http.DefaultClient)

	opts = append(opts, debugOption,loggerOption, httpOption)

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

```