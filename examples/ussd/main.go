package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
	"github.com/techcraftt/tigosdk"
	"github.com/techcraftt/tigosdk/ussd"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type disburseInfo struct {
	Msisdn string  `json:"msisdn"`
	Amount float64 `json:"amount"`
}

const (
	TIGO_USERNAME            = "TIGO_USERNAME"
	TIGO_PASSWORD            = "TIGO_PASSWORD"
	TIGO_PASSWORD_GRANT_TYPE = "TIGO_PASSWORD_GRANT_TYPE"
	TIGO_ACCOUNT_NAME        = "TIGO_ACCOUNT_NAME"
	TIGO_ACCOUNT_MSISDN      = "TIGO_ACCOUNT_MSISDN"
	TIGO_BRAND_ID            = "TIGO_BRAND_ID"
	TIGO_BILLER_CODE         = "TIGO_BILLER_CODE"
	TIGO_GET_TOKEN_URL       = "TIGO_GET_TOKEN_URL"
	TIGO_BILL_URL            = "TIGO_BILL_URL"
	TIGO_A2W_URL             = "TIGO_A2W_URL"
	TIGO_NAMECHECK_URL       = "TIGO_NAMECHECK_URL"
	TIGO_W2A_URL             = "TIGO_W2A_URL"
)

type App struct {
	Service ussd.Client
}

func (app *App) disburseHandler(writer http.ResponseWriter, request *http.Request) {
	var info disburseInfo
	//	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(request.Body).Decode(&info)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	req := ussd.AccountToWalletRequest{

		Type:        "REQMFCI",
		ReferenceID: refid,
		Msisdn:      app.Service.Conf.AccountMSISDN,
		PIN:         app.Service.Conf.Password,
		Msisdn1:     info.Msisdn,
		Amount:      info.Amount,
		SenderName:  app.Service.Conf.AccountName,
		Language1:   "EN",
		BrandID:     app.Service.Conf.BrandID,
	}

	resp, err := app.Service.AccountToWallet(context.TODO(), req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(resp)

}

func (app *App) namesHandler(writer http.ResponseWriter, request *http.Request) {

	resp, err := app.Service.QuerySubscriberName(context.TODO(), request)
	if err != nil {
		return
	}
	x, err := xml.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/xml")
	writer.Write(x)
}

func (app *App) transactionHandler(writer http.ResponseWriter, request *http.Request) {

	resp, err := app.Service.WalletToAccount(context.TODO(), request)

	if err != nil {
		return
	}
	x, err := xml.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/xml")
	writer.Write(x)
}

type User struct {
	Name   string `json:"name"`
	RefID  string `json:"ref_id"`
	Status int    `json:"status"`
}

var errNotFound = errors.New("not found")

func MakeHandler(client ussd.Client) http.Handler {

	app := App{client}

	router := mux.NewRouter()

	router.HandleFunc(app.Service.Conf.NameCheckRequestURL, app.namesHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(app.Service.Conf.WalletToAccountRequestURL, app.transactionHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc("/api/tigopesa/disburse", app.disburseHandler).Methods(http.MethodPost)

	return router
}

func loadFromEnv() (conf tigosdk.Configs, err error) {

	err = env.Load("tigo.env")
	conf = tigosdk.Configs{
		Username:                  os.Getenv(TIGO_USERNAME),
		Password:                  os.Getenv(TIGO_PASSWORD),
		PasswordGrantType:         os.Getenv(TIGO_PASSWORD_GRANT_TYPE),
		AccountName:               os.Getenv(TIGO_ACCOUNT_NAME),
		AccountMSISDN:             os.Getenv(TIGO_ACCOUNT_MSISDN),
		BrandID:                   os.Getenv(TIGO_BRAND_ID),
		BillerCode:                os.Getenv(TIGO_BILLER_CODE),
		GetTokenRequestURL:        os.Getenv(TIGO_GET_TOKEN_URL),
		PushPayBillRequestURL:     os.Getenv(TIGO_BILL_URL),
		AccountToWalletRequestURL: os.Getenv(TIGO_A2W_URL),
		WalletToAccountRequestURL: os.Getenv(TIGO_W2A_URL),
		NameCheckRequestURL:       os.Getenv(TIGO_NAMECHECK_URL),
	}

	return
}

func main() {
	conf, err := loadFromEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	usersMap := map[string]User{
		"12345678": {
			Name:   "Pius Alfred Shop",
			RefID:  "12345678",
			Status: 0,
		},
		"23456789": {
			Name:   "St. Jane School",
			RefID:  "23456789",
			Status: 1,
		},
		"34567890": {
			Name:   "Uhuru Stadium",
			RefID:  "34567890",
			Status: 2,
		},
		"22473478": {
			Name:   "Jamesson Club",
			RefID:  "22473478",
			Status: 2,
		},
	}

	check := checker{usersMap}

	c := ussd.NewClient(conf, nil, check.nameFunc, check.w2aFunc)

	handler := MakeHandler(*c)

	server := http.Server{
		Addr:              ":8090",
		Handler:           handler,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		return
	}

}

type checker struct {
	Users map[string]User
}

func (c *checker) w2aFunc(ctx context.Context, request ussd.WalletToAccountRequest) ussd.WalletToAccountResponse {

	user, found := c.checkUser(request.CustomerReferenceID)
	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	if ! found{
		resp := ussd.WalletToAccountResponse{

			Type:             tigosdk.SYNC_BILLPAY_RESPONSE,
			TxnID:            request.TxnID,
			RefID:            refid,
			Result:           "TF",
			ErrorCode:        "error010",
			ErrorDescription: "User Not Found",
			Msisdn:           request.Msisdn,
			Flag:             "N",
			Content:          request.SenderName,
		}

		return resp
	}else{
		if user.Status ==1 {
			resp := ussd.WalletToAccountResponse{
				Type:             tigosdk.SYNC_BILLPAY_RESPONSE,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        ussd.W2A_ERR_INVALID_CUSTOMER_REF_NUMBER,
				ErrorDescription: "Invalid Customer ref Number",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp
		}else if user.Status ==2{
			resp := ussd.WalletToAccountResponse{
				Type:             tigosdk.SYNC_BILLPAY_RESPONSE,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        ussd.W2A_ERR_CUSTOMER_REF_NUM_LOCKED,
				ErrorDescription: "Customer Locked",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp
		}
	}
	resp := ussd.WalletToAccountResponse{
		Type:             tigosdk.SYNC_BILLPAY_RESPONSE,
		TxnID:            request.TxnID,
		RefID:            refid,
		Result:           "TS",
		ErrorCode:        "error000",
		ErrorDescription: "Transaction Successful",
		Msisdn:           request.Msisdn,
		Flag:             "Y",
		Content:          request.SenderName,
	}
	return resp
}

func (c *checker) nameFunc(ctx context.Context, request ussd.SubscriberNameRequest) ussd.SubscriberNameResponse {

	user, found := c.checkUser(request.CustomerReferenceID)
	if !found {
		resp := ussd.SubscriberNameResponse{
			Type:      tigosdk.SYNC_LOOKUP_RESPONSE,
			Result:    "TF",
			ErrorCode: "error010",
			ErrorDesc: "Not found",
			Msisdn:    request.Msisdn,
			Flag:      "N",
			Content:   "User is not known",
		}

		return resp
	} else {
		if user.Status == 1 {

			resp := ussd.SubscriberNameResponse{
				Type:      tigosdk.SYNC_LOOKUP_RESPONSE,
				Result:    "TF",
				ErrorCode: "error020",
				ErrorDesc: "Transaction Failed: User Suspended",
				Msisdn:    request.Msisdn,
				Flag:      "N",
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp
		}

		if user.Status == 2 {
			resp := ussd.SubscriberNameResponse{
				Type:      tigosdk.SYNC_LOOKUP_RESPONSE,
				Result:    "TF",
				ErrorCode: "error030",
				ErrorDesc: "Transaction Failed: Format not known",
				Msisdn:    request.Msisdn,
				Flag:      "N",
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp
		}

		resp := ussd.SubscriberNameResponse{
			Type:      tigosdk.SYNC_LOOKUP_RESPONSE,
			Result:    "TS",
			ErrorCode: "error000",
			ErrorDesc: "Transaction Successfully",
			Msisdn:    request.Msisdn,
			Flag:      "Y",
			Content:   fmt.Sprintf("%s", user.Name),
		}
		return resp
	}

}

func (c *checker) checkUser(refid string) (User, bool) {
	fmt.Printf("checking %s\n", refid)

	user, found := c.Users[refid]

	return user, found
}
