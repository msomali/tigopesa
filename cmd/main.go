package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/techcraftt/tigosdk"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd"
	"io/ioutil"
	"net/http"
	"strconv"
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

	router.HandleFunc("/api/tigopesa/disburse", app.disburseHandler).Methods(http.MethodPost)

	return router
}

type disburseInfo struct {
	Msisdn string `json:"msisdn"`
	Amount float64 `json:"amount"`
}

type App struct {
	svc tigosdk.Service
	conf tigosdk.Configs
}

func (app App) transactionHandler(w http.ResponseWriter, request *http.Request) {

	var req ussd.W2ARequest
	xmlBody, err := ioutil.ReadAll(request.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = xml.Unmarshal(xmlBody, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err2 := app.svc.WalletToAccount(context.TODO(), req)
	if err2 != nil {
		return
	}

	x, err := xml.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write(x)
	if err != nil {
		return
	}

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

	resp, err := app.svc.QuerySubscriberName(context.TODO(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (app App) disburseHandler(w http.ResponseWriter, request *http.Request) {

	var info disburseInfo
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(request.Body).Decode(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	req := ussd.A2WRequest{

		Type:        "REQMFCI",
		ReferenceID: refid,
		Msisdn:      app.conf.AccountMSISDN,
		PIN:         app.conf.Password,
		Msisdn1:     info.Msisdn,
		Amount:      info.Amount,
		SenderName:  "Bethuel Charles",
		Language1:   "EN",
		BrandID:     app.conf.BrandID,
	}


	xmlstring, err := xml.MarshalIndent(req, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	xmlstring = []byte(xml.Header + string(xmlstring))

	r, err := http.NewRequest("POST", app.conf.A2WReqURL, bytes.NewBuffer(xmlstring)) // URL-encoded payload
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	r.Header.Add("Content-Type", "application/xml")
	
	client := http.Client{
		Timeout: time.Minute,
	}

	res, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Println(res.Status)
	defer res.Body.Close()
	xmlBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var xmlResponse ussd.A2WResponse

	err = xml.Unmarshal(xmlBody, &xmlResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	x, err := xml.MarshalIndent(xmlResponse, "", "  ")
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
