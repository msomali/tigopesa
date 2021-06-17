package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	TIGO_USERNAME      = "TIGO_USERNAME"
	TIGO_PASSWORD      = "TIGO_PASSWORD"
	TIGO_BILLER_MSISDN = "TIGO_BILLER_MSISDN"
	TIGO_BILLER_CODE   = "TIGO_BILLER_CODE"
	TIGO_GET_TOKEN_URL = "TIGO_GET_TOKEN_URL"
	TIGO_BILL_URL      = "TIGO_BILL_URL"
	TIGO_BASE_URL      = "TIGO_BASE_URL"
)

type (
	pushpayInitiatorRequest struct {
		CustomerMSSID int64  `json:"customer"`
		Amount        int    `json:"amount"`
		Remarks       string `json:"remarks, omitempty"`
	}

	app struct {
		pushpay     push.Service
		transaction map[string]push.BillPayRequest
	}
)

func main() {
	config, err := loadFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}

	a := &app{
		pushpay: push.NewClient(&tigo.BaseClient{
			Config:     config,
			HttpClient: nil,
			Ctx:        nil,
			Timeout:    0,
			Logger:     nil,
		}, func(ctx context.Context, request push.BillPayCallbackRequest) *push.BillPayResponse {
			return nil
		}),
		transaction: map[string]push.BillPayRequest{},
	}

	router := mux.NewRouter()

	router.HandleFunc("/tigopesa/pushpay", a.pushPayHandler).Methods(http.MethodPost)
	router.HandleFunc("/tigopesa/pushpay/callback", a.pushPayCallbackHandler).Methods(http.MethodPost)

	s := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err.Error())
	}
}

func (a *app) pushPayHandler(w http.ResponseWriter, r *http.Request) {
	var req pushpayInitiatorRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	billRequest := push.BillPayRequest{
		CustomerMSISDN: strconv.FormatInt(req.CustomerMSSID, 10),
		Amount:         req.Amount,
		Remarks:        req.Remarks,
		ReferenceID:    fmt.Sprintf("%s%d", os.Getenv(TIGO_BILLER_CODE), time.Now().Local().Unix()),
	}
	billPayResponse, err := a.pushpay.BillPay(context.Background(), billRequest)
	if err != nil {
		log.Printf("PushBillPay request failed error: %s", err.Error())
		return
	}

	// keep record of successfully transaction if status was success initiated.
	if billPayResponse.ResponseStatus {
		a.transaction[billRequest.ReferenceID] = billRequest
	}

	return
}

func (a *app) pushPayCallbackHandler(w http.ResponseWriter, r *http.Request) {
	a.pushpay.BillPayCallback(context.Background())
}

func (a *app) callbackProvider(ctx context.Context, billPayRequest push.BillPayCallbackRequest) *push.BillPayResponse {
	if !billPayRequest.Status {
		return &push.BillPayResponse{
			ResponseCode:        "BILLER-18-3020-E",
			ResponseStatus:      false,
			ResponseDescription: "Callback failed",
			ReferenceID:         billPayRequest.ReferenceID,
		}
	}

	return &push.BillPayResponse{
		ResponseCode:        "BILLER-18-0000-S",
		ResponseStatus:      true,
		ResponseDescription: "Callback successful",
		ReferenceID:         billPayRequest.ReferenceID,
	}
}

func loadFromEnv() (conf tigo.Config, err error) {
	var billerMSISDN int64

	err = env.Load("tigo.env")
	if err != nil {
		log.Fatalln(err.Error())
	}

	billerMSISDN, err = strconv.ParseInt(os.Getenv(TIGO_BILLER_MSISDN), 10, 64)

	conf = tigo.Config{
		Username:              os.Getenv(TIGO_USERNAME),
		Password:              os.Getenv(TIGO_PASSWORD),
		BillerCode:            os.Getenv(TIGO_BILLER_CODE),
		BillerMSISDN:          strconv.FormatInt(billerMSISDN, 10),
		GetTokenRequestURL:    os.Getenv(TIGO_GET_TOKEN_URL),
		PushPayBillRequestURL: os.Getenv(TIGO_BILL_URL),
		ApiBaseURL:            os.Getenv(TIGO_BASE_URL),
	}

	return
}
