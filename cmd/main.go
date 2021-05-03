package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/techcraftt/tigosdk"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd"
	"io/ioutil"
	"net/http"
	"time"
)
import "github.com/gorilla/mux"

func MakeHandler(svc tigosdk.Service) http.Handler {

	app := App{svc: svc}
	router := mux.NewRouter()

	//namecheck
	//
	//http://192.168.176.244:8090/api/tigopesa/c2b/users/names
	//
	//payment
	//
	//http://192.168.176.244:8090/api/tigopesa/c2b/transactions

	router.HandleFunc("/api/tigopesa/c2b/users/names", app.namesHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc("/api/tigopesa/c2b/transactions", app.transactionHandler).Methods(http.MethodPost, http.MethodGet)

	return router
}

type App struct {
	svc tigosdk.Service
}

func (app App) transactionHandler(writer http.ResponseWriter, request *http.Request) {

}

func (app App) namesHandler(w http.ResponseWriter, request *http.Request) {

	var req ussd.SubscriberNameRequest

	xmlBody, err := ioutil.ReadAll(request.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = xml.Unmarshal(xmlBody, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err2 := app.svc.QuerySubscriberName(context.TODO(), req)
	if err2 != nil {
		return
	}

	x, err := xml.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func main() {
	conf := tigosdk.Configs{
		Username:          "",
		Password:          "",
		PasswordGrantType: "",
		AccountName:       "",
		AccountMSISDN:     "",
		BrandID:           "",
		BillerCode:        "",
		GetTokenURL:       "",
		BillURL:           "",
		A2WReqURL:         "",
	}

	namechecker := func(ctx context.Context, request ussd.NameCheckRequest) (ussd.NameCheckResponse, error) {
		resp := ussd.NameCheckResponse{
			Result:    "TS",
			ErrorCode: "error000",
			ErrorDesc: "Transaction Successfully So no Err",
			Msisdn:    "255712915790",
			Flag:      "Y",
			Content:   "this is content",
		}
		return resp, nil
	}

	callbacker := func(ctx context.Context, request push.BillPayCallbackRequest) {
		fmt.Printf("do nothing")
	}

	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	svc := tigosdk.NewTigoClient(&httpClient, namechecker, callbacker, conf)

	handler := MakeHandler(svc)

	server := http.Server{
		Addr:              ":8090",
		Handler:           handler,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
