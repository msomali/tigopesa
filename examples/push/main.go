package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
)

type (
	pushpayInitiatorRequest struct {
		CustomerMSSID int64  `json:"customer"`
		Amount        int    `json:"amount"`
		Remarks       string `json:"remarks,omitempty"`
	}

	app struct {
		push *push.Client
	}
)

func pushPayCallbackHandler() push.CallbackHandlerFunc {
	return func(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {
		if !request.Status {
			return push.CallbackResponse{
				ResponseCode:        "BILLER-18-3020-E",
				ResponseStatus:      false,
				ResponseDescription: "Callback failed",
				ReferenceID:         request.ReferenceID,
			}, nil
		}

		return push.CallbackResponse{
			ResponseCode:        "BILLER-18-0000-S",
			ResponseStatus:      true,
			ResponseDescription: "Callback successful",
			ReferenceID:         request.ReferenceID,
		}, nil

	}
}

func pushHandler()func(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error){
	return func(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {
		return push.CallbackResponse{},nil
	}
}

func main() {
	config, err := loadFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}

	sdk := tigo.NewBaseClient()

	pushClient := &push.Client{
		BaseClient:      sdk,
		Config:          config,
		CallbackHandler: pushPayCallbackHandler(),
	}

	a := &app{
		push: pushClient,
	}

	router := mux.NewRouter()

	router.HandleFunc("/tigopesa/pushpay", a.pushPayHandler).Methods(http.MethodPost)
	router.HandleFunc("/tigopesa/pushpay/callback", a.push.Callback).Methods(http.MethodPost)

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

	billRequest := push.PayRequest{
		CustomerMSISDN: strconv.FormatInt(req.CustomerMSSID, 10),
		Amount:         req.Amount,
		Remarks:        req.Remarks,
		ReferenceID:    fmt.Sprintf("%s%d", os.Getenv("TIGO_BILLER_CODE"), time.Now().Local().Unix()),
	}

	response, err := a.push.Pay(context.Background(), billRequest)
	if err != nil {
		log.Printf("PushBillPay request failed error: %s", err.Error())
		return
	}

	fmt.Printf("%v\n", response)
	return
}

func loadFromEnv() (conf *push.Config, err error) {
	var billerMSISDN int64

	err = godotenv.Load("tigo.env")
	if err != nil {
		log.Fatalln(err.Error())
	}

	billerMSISDN, err = strconv.ParseInt(os.Getenv("TIGO_BILLER_MSISDN"), 10, 64)

	conf = &push.Config{
		Username:     os.Getenv("TIGO_USERNAME"),
		Password:     os.Getenv("TIGO_PASSWORD"),
		BillerCode:   os.Getenv("TIGO_BILLER_CODE"),
		BillerMSISDN: strconv.FormatInt(billerMSISDN, 10),
		GetTokenURL:  os.Getenv("TIGO_GET_TOKEN_URL"),
		PushPayURL:   os.Getenv("TIGO_BILL_URL"),
		ApiBaseURL:   os.Getenv("TIGO_BASE_URL"),
	}

	return
}
