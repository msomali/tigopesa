package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
	"github.com/techcraftt/tigosdk"
	"github.com/techcraftt/tigosdk/aw"
	"github.com/techcraftt/tigosdk/pkg/conf"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/wa"
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
	TIGO_USERNAME         = "TIGO_USERNAME"
	TIGO_PASSWORD         = "TIGO_PASSWORD"
	TIGO_ACCOUNT_NAME     = "TIGO_ACCOUNT_NAME"
	TIGO_ACCOUNT_MSISDN   = "TIGO_ACCOUNT_MSISDN"
	TIGO_BRAND_ID         = "TIGO_BRAND_ID"
	TIGO_BILLER_CODE      = "TIGO_BILLER_CODE"
	TIGO_A2W_URL          = "TIGO_A2W_URL"
	TIGO_NAMECHECK_URL    = "TIGO_NAMECHECK_URL"
	TIGO_W2A_URL          = "TIGO_W2A_URL"
	TIGO_DISBURSEMENT_PIN = "TIGO_DISBURSEMENT_PIN"
)

type App struct {
	wa *wa.Client
	aw *aw.Client
}

func (app *App) disburseHandler(writer http.ResponseWriter, request *http.Request) {
	var info disburseInfo
	//	// Try to decode the request body into the struct. If there is an error,
	// respond to the pkg with the error message and a 400 status code.
	err := json.NewDecoder(request.Body).Decode(&info)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	refid := fmt.Sprintf("PCT%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	req := aw.DisburseRequest{
		Type:        tigosdk.REQMFCI,
		ReferenceID: refid,
		Msisdn:      app.aw.Config.AccountMSISDN,
		PIN:         app.aw.Config.PIN,
		Msisdn1:     info.Msisdn,
		Amount:      info.Amount,
		SenderName:  app.aw.Config.AccountName,
		Language1:   "EN",
		BrandID:     app.aw.Config.BrandID,
	}

	resp, err := app.aw.Disburse(context.TODO(), req)
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

func MakeHandler(client1 *wa.Client, client2 *aw.Client) http.Handler {

	app := App{
		wa: client1,
		aw: client2,
	}

	router := mux.NewRouter()

	router.HandleFunc(app.wa.NamecheckURL, app.wa.HandleNameQuery).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(app.wa.RequestURL, app.wa.HandlePayment).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(app.aw.RequestURL, app.disburseHandler).Methods(http.MethodPost)

	return router
}

func loadFromEnv() (config *conf.Config, err error) {

	err = env.Load("tigo.env")

	config = &conf.Config{
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

	if err != nil {
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

	keeper := checker{usersMap}

	var opts []tigo.ClientOption

	opts = append(opts,
		tigo.WithContext(context.Background()),
		tigo.WithTimeout(time.Second*30),
		tigo.WithLogger(os.Stderr),
		tigo.WithHTTPClient(http.DefaultClient),
	)

	_, pay, disburse := conf.Split()

	bc := tigo.NewBaseClient(tigo.WithDebugMode(true))

	waClient := &wa.Client{
		BaseClient:       bc,
		Config:           pay,
		PaymentHandler:   keeper,
		NameQueryHandler: keeper,
	}

	awClient := &aw.Client{
		Config:     disburse,
		BaseClient: bc,
	}

	handler := MakeHandler(waClient, awClient)

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

func (c checker) HandlePaymentRequest(ctx context.Context, request wa.PayRequest) (wa.PayResponse, error) {

	user, found := c.checkUser(request.CustomerReferenceID)
	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	if !found {
		resp := wa.PayResponse{

			Type:             tigosdk.SyncBillPayResponse,
			TxnID:            request.TxnID,
			RefID:            refid,
			Result:           tigosdk.FailedTxnResult,
			ErrorCode:        tigosdk.ErrInvalidCustomerRefNumber,
			ErrorDescription: "User Not Found",
			Msisdn:           request.Msisdn,
			Flag:             tigosdk.NoFlag,
			Content:          request.SenderName,
		}

		return resp, nil
	} else {
		if user.Status == 1 {
			resp := wa.PayResponse{
				Type:             tigosdk.SyncBillPayResponse,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        tigosdk.ErrInvalidCustomerRefNumber,
				ErrorDescription: "Invalid Customer ref Number",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp, nil
		} else if user.Status == 2 {
			resp := wa.PayResponse{
				Type:             tigosdk.SyncBillPayResponse,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        tigosdk.ErrCustomerRefNumLocked,
				ErrorDescription: "Customer Locked",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp, nil
		}
	}
	resp := wa.PayResponse{
		Type:             tigosdk.SyncBillPayResponse,
		TxnID:            request.TxnID,
		RefID:            refid,
		Result:           "TS",
		ErrorCode:        "error000",
		ErrorDescription: "Transaction Successful",
		Msisdn:           request.Msisdn,
		Flag:             "Y",
		Content:          request.SenderName,
	}
	return resp, nil
}

func (c checker) HandleSubscriberNameQuery(ctx context.Context, request wa.NameRequest) (wa.NameResponse, error) {

	user, found := c.checkUser(request.CustomerReferenceID)
	if !found {
		resp := wa.NameResponse{
			Type:      tigosdk.SyncLookupResponse,
			Result:    "TF",
			ErrorCode: "error010",
			ErrorDesc: "Not found",
			Msisdn:    request.Msisdn,
			Flag:      "N",
			Content:   "User is not known",
		}

		return resp, nil
	} else {
		if user.Status == 1 {

			resp := wa.NameResponse{
				Type:      tigosdk.SyncLookupResponse,
				Result:    "TF",
				ErrorCode: "error020",
				ErrorDesc: "Transaction Failed: User Suspended",
				Msisdn:    request.Msisdn,
				Flag:      "N",
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp, nil
		}

		if user.Status == 2 {
			resp := wa.NameResponse{
				Type:      tigosdk.SyncLookupResponse,
				Result:    "TF",
				ErrorCode: tigosdk.ErrNameInvalidFormat,
				ErrorDesc: "Transaction Failed: Format not known",
				Msisdn:    request.Msisdn,
				Flag:      "N",
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp, nil
		}

		resp := wa.NameResponse{
			Type:      tigosdk.SyncLookupResponse,
			Result:    tigosdk.SucceededTxnResult,
			ErrorCode: tigosdk.NoNamecheckErr,
			ErrorDesc: "Transaction Successfully",
			Msisdn:    request.Msisdn,
			Flag:      tigosdk.YesFlag,
			Content:   fmt.Sprintf("%s", user.Name),
		}
		return resp, nil
	}

}

func (c *checker) checkUser(refid string) (User, bool) {
	fmt.Printf("checking %s\n", refid)

	user, found := c.Users[refid]

	return user, found
}
