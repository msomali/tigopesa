package main
//
//import (
//	"context"
//	"encoding/json"
//	"encoding/xml"
//	"errors"
//	"fmt"
//	"github.com/techcraftt/tigosdk"
//	"github.com/techcraftt/tigosdk/push"
//	"github.com/techcraftt/tigosdk/ussd"
//	"io/ioutil"
//	"net/http"
//	"strconv"
//	"time"
//)
//import "github.com/gorilla/mux"
//
//type User struct {
//	Name string `json:"name"`
//	RefID string `json:"ref_id"`
//	Status int `json:"status"`
//}
//
//var errNotFound = errors.New("not found")
//
//func checkUser(refid string)(User,bool){
//	fmt.Printf("checking %s\n",refid)
//
//	usersMap := map[string]User{
//		"12345678": {
//			Name:   "Pius Alfred Shop",
//			RefID:  "12345678",
//			Status: 0,
//		},
//		"23456789":{
//			Name:   "St. Jane School",
//			RefID:  "23456789",
//			Status: 1,
//		},
//		"34567890":{
//			Name:   "Uhuru Stadium",
//			RefID:  "34567890",
//			Status: 2,
//		},
//		"22473478":{
//			Name:   "Jamesson Club",
//			RefID:  "22473478",
//			Status: 2,
//		},
//	}
//
//	user, found := usersMap[refid]
//
//	return user, found
//}
//
//func MakeHandler(svc tigosdk.Service) http.Handler {
//
//	app := App{svc: svc}
//	router := mux.NewRouter()
//
//	//namecheck
//	//
//	//http://192.168.176.244:8090/api/tigopesa/c2b/users/names
//	//
//	//payment
//	//
//	//http://192.168.176.244:8090/api/tigopesa/c2b/transactions
//
//	router.HandleFunc("/api/tigopesa/c2b/users/names", app.namesHandler).Methods(http.MethodPost, http.MethodGet)
//
//	router.HandleFunc("/api/tigopesa/c2b/transactions", app.transactionHandler).Methods(http.MethodPost, http.MethodGet)
//
//	router.HandleFunc("/api/tigopesa/disburse", app.disburseHandler).Methods(http.MethodPost)
//
//	return router
//}
//
//type disburseInfo struct {
//	Msisdn string `json:"msisdn"`
//	Amount float64 `json:"amount"`
//}
//
//type App struct {
//	svc tigosdk.Service
//	conf tigosdk.Configs
//}
//
//func (app App) transactionHandler(w http.ResponseWriter, request *http.Request) {
//
//	var req ussd.WalletToAccountRequest
//	xmlBody, err := ioutil.ReadAll(request.Body)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err = xml.Unmarshal(xmlBody, &req)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	resp, err2 := app.svc.WalletToAccount(context.TODO(), req)
//	if err2 != nil {
//		return
//	}
//
//	x, err := xml.MarshalIndent(resp, "", "  ")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/xml")
//	_, err = w.Write(x)
//	if err != nil {
//		return
//	}
//
//}
//
//func (app App) namesHandler(w http.ResponseWriter, request *http.Request) {
//
//	var req ussd.SubscriberNameRequest
//
//	xmlBody, err := ioutil.ReadAll(request.Body)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Try to decode the request body into the struct. If there is an error,
//	// respond to the client with the error message and a 400 status code.
//	err = xml.Unmarshal(xmlBody, &req)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	resp, err := app.svc.QuerySubscriberName(context.TODO(), req)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	x, err := xml.MarshalIndent(resp, "", "  ")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/xml")
//	w.Write(x)
//}
//
//func (app App) disburseHandler(w http.ResponseWriter, request *http.Request) {
//
//	var info disburseInfo
//	// Try to decode the request body into the struct. If there is an error,
//	// respond to the client with the error message and a 400 status code.
//	err := json.NewDecoder(request.Body).Decode(&info)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))
//
//	req := ussd.AccountToWalletRequest{
//
//		Type:        "REQMFCI",
//		ReferenceID: refid,
//		Msisdn:      app.conf.AccountMSISDN,
//		PIN:         app.conf.Password,
//		Msisdn1:     info.Msisdn,
//		Amount:      info.Amount,
//		SenderName:  "Bethuel Charles",
//		Language1:   "EN",
//		BrandID:     app.conf.BrandID,
//	}
//
//	resp, err := app.svc.AccountToWallet(context.TODO(),req)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//
//	json.NewEncoder(w).Encode(resp)
//}
//
//func main() {
//	conf := tigosdk.Configs{
//		Username:                  "",
//		Password:                  "",
//		PasswordGrantType:         "",
//		AccountName:               "",
//		AccountMSISDN:             "",
//		BrandID:                   "",
//		BillerCode:                "",
//		GetTokenRequestURL:        "",
//		PushPayBillRequestURL:     "",
//		AccountToWalletRequestURL: "",
//	}
//
//	namechecker := func(ctx context.Context, request ussd.SubscriberNameRequest) (ussd.NameCheckResponse, error) {
//
//		user, found := checkUser(request.CustomerReferenceID)
//
//		if !found{
//			resp := ussd.NameCheckResponse{
//				Result:    "TF",
//				ErrorCode: "error010",
//				ErrorDesc: "Not found",
//				Msisdn:    request.Msisdn,
//				Flag:      "N",
//				Content:   "User is not known",
//			}
//
//			return resp, nil
//		}else {
//			if user.Status == 1{
//
//				resp := ussd.NameCheckResponse{
//					Result:    "TF",
//					ErrorCode: "error020",
//					ErrorDesc: "Transaction Failed: User Suspended",
//					Msisdn:    request.Msisdn,
//					Flag:      "N",
//					Content:   fmt.Sprintf("user name %s", user.Name),
//				}
//
//				return resp, nil
//			}
//
//			if user.Status ==2{
//				resp := ussd.NameCheckResponse{
//					Result:    "TF",
//					ErrorCode: "error030",
//					ErrorDesc: "Transaction Failed: Format not known",
//					Msisdn:    request.Msisdn,
//					Flag:      "N",
//					Content:   fmt.Sprintf("user name %s", user.Name),
//				}
//
//				return resp, nil
//			}
//
//			resp := ussd.NameCheckResponse{
//				Result:    "TS",
//				ErrorCode: "error000",
//				ErrorDesc: "Transaction Successfully",
//				Msisdn:    request.Msisdn,
//				Flag:      "Y",
//				Content:   "transaction successful, user known and valid",
//			}
//			return resp, nil
//		}
//
//	}
//
//	callbacker := func(ctx context.Context, request push.BillPayCallbackRequest) {
//		fmt.Printf("do nothing")
//	}
//
//	httpClient := http.Client{
//		Timeout: 60 * time.Second,
//	}
//
//	svc := tigosdk.NewTigoClient(&httpClient, namechecker, callbacker, conf)
//
//	handler := MakeHandler(svc)
//
//	server := http.Server{
//		Addr:              ":8090",
//		Handler:           handler,
//		ReadTimeout:       30 * time.Second,
//		ReadHeaderTimeout: 30 * time.Second,
//		WriteTimeout:      30 * time.Second,
//	}
//
//	err := server.ListenAndServe()
//	if err != nil {
//		return
//	}
//}
