package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
	"github.com/techcraftt/tigosdk"
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
		Remarks       string `json:"remarks"`
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
		pushpay: push.NewClientFromConfig(config),
	}

	r := mux.NewRouter()

	r.HandleFunc("/tigopesa/pushpay", a.pushPayHandler).Methods(http.MethodPost)
	r.HandleFunc("/tigopesa/pushpay/callback", a.pushPayCallbackHandler).Methods(http.MethodPost)

	s := http.Server{
		Addr:              ":8090",
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
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
		CustomerMSISDN: req.CustomerMSSID,
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
	var req push.BillPayCallbackRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// callback failed if we failed to decode the request sent from tigo
		json.NewEncoder(w).Encode(&push.BillPayResponse{
			ResponseCode:        "BILLER-18-3020-E",
			ResponseStatus:      false,
			ResponseDescription: "Callback failed",
			ReferenceID:         req.ReferenceID,
		})

		return
	}

	//check if callback status is successful and the request originated from our API.
	if req.Status {
		// verify if received billpay reference initiated by our api.
		trans, ok := a.transaction[req.ReferenceID]
		if ok && trans.Amount == req.Amount && trans.ReferenceID == req.ReferenceID {
			json.NewEncoder(w).Encode(&push.BillPayResponse{
				ResponseCode:        "BILLER-18-0000-S",
				ResponseStatus:      true,
				ResponseDescription: "Callback successful",
				ReferenceID:         req.ReferenceID,
			})
		}
	}

	// respond with failed callback
	json.NewEncoder(w).Encode(&push.BillPayResponse{
		ResponseCode:        "BILLER-18-3020-E",
		ResponseStatus:      false,
		ResponseDescription: "Callback failed",
		ReferenceID:         req.ReferenceID,
	})

	return
}

func loadFromEnv() (conf tigosdk.Config, err error) {
	var billerMSISDN int64

	err = env.Load("tigo.env")

	billerMSISDN, err = strconv.ParseInt(os.Getenv(TIGO_BILLER_MSISDN), 10, 64)

	conf = tigosdk.Config{
		Username:              os.Getenv(TIGO_USERNAME),
		Password:              os.Getenv(TIGO_PASSWORD),
		BillerCode:            os.Getenv(TIGO_BILLER_CODE),
		BillerMSISDN:          billerMSISDN,
		GetTokenRequestURL:    os.Getenv(TIGO_GET_TOKEN_URL),
		PushPayBillRequestURL: os.Getenv(TIGO_BILL_URL),
		ApiBaseURL:            os.Getenv(TIGO_BASE_URL),
	}

	return
}
