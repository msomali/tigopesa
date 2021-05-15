package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
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
	TIGO_ACCOUNT_NAME        = "TIGO_ACCOUNT_NAME"
	TIGO_ACCOUNT_MSISDN      = "TIGO_ACCOUNT_MSISDN"
	TIGO_BRAND_ID            = "TIGO_BRAND_ID"
	TIGO_BILLER_CODE         = "TIGO_BILLER_CODE"
	TIGO_A2W_URL             = "TIGO_A2W_URL"
	TIGO_NAMECHECK_URL       = "TIGO_NAMECHECK_URL"
	TIGO_W2A_URL             = "TIGO_W2A_URL"
	TIGO_DISBURSEMENT_PIN    = "TIGO_DISBURSEMENT_PIN"
)

type App struct {
	USSDClient *ussd.Client
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

	refid := fmt.Sprintf("PCT%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	req := ussd.AccountToWalletRequest{

		Type:        ussd.REQMFCI,
		ReferenceID: refid,
		Msisdn:      app.USSDClient.Config.AccountMSISDN,
		PIN:         app.USSDClient.Config.AccountToWalletRequestPIN,
		Msisdn1:     info.Msisdn,
		Amount:      info.Amount,
		SenderName:  app.USSDClient.Config.AccountName,
		Language1:   "EN",
		BrandID:     app.USSDClient.Config.BrandID,
	}

	//logger := log.New(os.Stdout,"disburse",1)
	//logger.Printf("disburse request %v\n",req)

	resp, err := app.USSDClient.AccountToWalletHandler(context.TODO(), req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	//logger.Printf("disburse response %v\n",resp)

	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(resp)

}


type User struct {
	Name   string `json:"name"`
	RefID  string `json:"ref_id"`
	Status int    `json:"status"`
}

func MakeHandler(client *ussd.Client) http.Handler {

	app := App{client}

	router := mux.NewRouter()

	router.HandleFunc(app.USSDClient.NameCheckRequestURL, app.USSDClient.SubscriberNameHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(app.USSDClient.WalletToAccountRequestURL, app.USSDClient.WalletToAccountHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(app.USSDClient.AccountToWalletRequestURL, app.disburseHandler).Methods(http.MethodPost)

	return router
}

func loadFromEnv() (conf ussd.Config, err error) {

	err = env.Load("tigo.env")
	conf = ussd.Config{
		Username:                  os.Getenv(TIGO_USERNAME),
		Password:                  os.Getenv(TIGO_PASSWORD),
		AccountName:               os.Getenv(TIGO_ACCOUNT_NAME),
		AccountMSISDN:             os.Getenv(TIGO_ACCOUNT_MSISDN),
		BrandID:                   os.Getenv(TIGO_BRAND_ID),
		BillerCode:                os.Getenv(TIGO_BILLER_CODE),
		AccountToWalletRequestURL: os.Getenv(TIGO_A2W_URL),
		WalletToAccountRequestURL: os.Getenv(TIGO_W2A_URL),
		AccountToWalletRequestPIN: os.Getenv(TIGO_DISBURSEMENT_PIN),
		NameCheckRequestURL:       os.Getenv(TIGO_NAMECHECK_URL),
	}

	if err != nil{
		panic(err)
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

	var opts [] ussd.ClientOption

	opts = append(opts,
		ussd.WithContext(context.Background()),
		ussd.WithTimeout(time.Second*30),
		ussd.WithLogger(os.Stderr),
		ussd.WithHTTPClient(http.DefaultClient),
		)

	c := ussd.NewClient(conf,check.w2aFunc,check.nameFunc,opts...)

	handler := MakeHandler(c)

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

func (c *checker) w2aFunc(ctx context.Context, request ussd.WalletToAccountRequest) (ussd.WalletToAccountResponse,error) {

	user, found := c.checkUser(request.CustomerReferenceID)
	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	if ! found{
		resp := ussd.WalletToAccountResponse{

			Type:             ussd.SyncBillPayResponse,
			TxnID:            request.TxnID,
			RefID:            refid,
			Result:           "TF",
			ErrorCode:        "error010",
			ErrorDescription: "User Not Found",
			Msisdn:           request.Msisdn,
			Flag:             "N",
			Content:          request.SenderName,
		}

		return resp,nil
	}else{
		if user.Status ==1 {
			resp := ussd.WalletToAccountResponse{
				Type:             ussd.SyncBillPayResponse,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        ussd.W2A_ERR_INVALID_CUSTOMER_REF_NUMBER,
				ErrorDescription: "Invalid Customer ref Number",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp,nil
		}else if user.Status ==2{
			resp := ussd.WalletToAccountResponse{
				Type:             ussd.SyncBillPayResponse,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        ussd.W2A_ERR_CUSTOMER_REF_NUM_LOCKED,
				ErrorDescription: "Customer Locked",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp,nil
		}
	}
	resp := ussd.WalletToAccountResponse{
		Type:             ussd.SyncBillPayResponse,
		TxnID:            request.TxnID,
		RefID:            refid,
		Result:           "TS",
		ErrorCode:        "error000",
		ErrorDescription: "Transaction Successful",
		Msisdn:           request.Msisdn,
		Flag:             "Y",
		Content:          request.SenderName,
	}
	return resp,nil
}

func (c *checker) nameFunc(ctx context.Context, request ussd.SubscriberNameRequest) (ussd.SubscriberNameResponse,error){

	user, found := c.checkUser(request.CustomerReferenceID)
	if !found {
		resp := ussd.SubscriberNameResponse{
			Type:      ussd.SyncLookupResponse,
			Result:    "TF",
			ErrorCode: "error010",
			ErrorDesc: "Not found",
			Msisdn:    request.Msisdn,
			Flag:      "N",
			Content:   "User is not known",
		}

		return resp,nil
	} else {
		if user.Status == 1 {

			resp := ussd.SubscriberNameResponse{
				Type:      ussd.SyncLookupResponse,
				Result:    "TF",
				ErrorCode: "error020",
				ErrorDesc: "Transaction Failed: User Suspended",
				Msisdn:    request.Msisdn,
				Flag:      "N",
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp,nil
		}

		if user.Status == 2 {
			resp := ussd.SubscriberNameResponse{
				Type:      ussd.SyncLookupResponse,
				Result:    "TF",
				ErrorCode: "error030",
				ErrorDesc: "Transaction Failed: Format not known",
				Msisdn:    request.Msisdn,
				Flag:      "N",
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp,nil
		}

		resp := ussd.SubscriberNameResponse{
			Type:      ussd.SyncLookupResponse,
			Result:    "TS",
			ErrorCode: "error000",
			ErrorDesc: "Transaction Successfully",
			Msisdn:    request.Msisdn,
			Flag:      "Y",
			Content:   fmt.Sprintf("%s", user.Name),
		}
		return resp,nil
	}

}

func (c *checker) checkUser(refid string) (User, bool) {
	fmt.Printf("checking %s\n", refid)

	user, found := c.Users[refid]

	return user, found
}
