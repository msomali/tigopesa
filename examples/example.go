package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/techcraftt/tigosdk"
	"github.com/techcraftt/tigosdk/aw"
	"github.com/techcraftt/tigosdk/pkg/conf"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/wa"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var Config = &conf.Config{
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

func makeApp() *App {

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

	bc := &tigo.BaseClient{
		HttpClient: http.DefaultClient,
		Ctx:        context.Background(),
		Timeout:    60 * time.Second,
		Logger:     os.Stdout,
		DebugMode:  true,
	}

	pc, wc, ac := Config.Split()

	p := &push.Client{
		Config:          pc,
		BaseClient:      bc,
		CallbackHandler: pushPayCallbackHandler(),
	}

	a := &aw.Client{
		Config:     ac,
		BaseClient: bc,
	}

	w := &wa.Client{
		BaseClient:       bc,
		Config:           wc,
		PaymentHandler:   keeper,
		NameQueryHandler: keeper,
	}

	app := &App{
		Config:   Config,
		push:     p,
		disburse: a,
		ussd:     w,
	}

	return app
}

type (
	disburseInfo struct {
		Msisdn string  `json:"msisdn"`
		Amount float64 `json:"amount"`
	}

	App struct {
		Config   *conf.Config
		push     *push.Client
		disburse *aw.Client
		ussd     *wa.Client
	}

	pushpayInitiatorRequest struct {
		CustomerMSSID int64  `json:"customer"`
		Amount        int    `json:"amount"`
		Remarks       string `json:"remarks,omitempty"`
	}
)

func (a *App) disburseHandler(writer http.ResponseWriter, request *http.Request) {
	var info disburseInfo
	//	// Try to decode the request body into the struct. If there is an error,
	// respond to the pkg with the error message and a 400 status code.
	err := json.NewDecoder(request.Body).Decode(&info)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	refid := fmt.Sprintf("PCT%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	resp, err := a.disburse.Disburse(context.TODO(), refid, info.Msisdn, info.Amount)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	//logger.Printf("disburse response %v\n",resp)

	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(resp)

}

func (a *App) makeHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc(a.ussd.NamecheckURL, a.ussd.HandleNameQuery).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc(a.ussd.RequestURL, a.ussd.HandlePayment).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/api/tigopesa/disburse", a.disburseHandler).Methods(http.MethodPost)
	router.HandleFunc("/tigopesa/pushpay", a.pushPayHandler).Methods(http.MethodPost)
	router.HandleFunc("/tigopesa/pushpay/callback", a.push.Callback).Methods(http.MethodPost)

	return router
}

func pushPayCallbackHandler() push.CallbackHandlerFunc {
	return func(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {
		if !request.Status {
			return push.CallbackResponse{
				ResponseCode:        push.FailureCode,
				ResponseStatus:      false,
				ResponseDescription: "Callback failed",
				ReferenceID:         request.ReferenceID,
			}, nil
		}

		return push.CallbackResponse{
			ResponseCode:        push.SuccessCode,
			ResponseStatus:      true,
			ResponseDescription: "Callback successful",
			ReferenceID:         request.ReferenceID,
		}, nil

	}
}

func (a *App) pushPayHandler(w http.ResponseWriter, r *http.Request) {
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

type User struct {
	Name   string `json:"name"`
	RefID  string `json:"ref_id"`
	Status int    `json:"status"`
}

type checker struct {
	Users map[string]User
}

func (c checker) HandlePaymentRequest(ctx context.Context, request wa.PayRequest) (wa.PayResponse, error) {

	user, found := c.checkUser(request.CustomerReferenceID)
	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	if !found {
		resp := wa.PayResponse{

			Type:             wa.SyncBillPayResponse,
			TxnID:            request.TxnID,
			RefID:            refid,
			Result:           tigosdk.FailedTxnResult,
			ErrorCode:        aw.ErrInvalidCustomerRefNumber,
			ErrorDescription: "User Not Found",
			Msisdn:           request.Msisdn,
			Flag:             tigosdk.NoFlag,
			Content:          request.SenderName,
		}

		return resp, nil
	} else {
		if user.Status == 1 {
			resp := wa.PayResponse{
				Type:             wa.SyncBillPayResponse,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        wa.ErrInvalidCustomerRefNumber,
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
				ErrorCode:        wa.ErrCustomerRefNumLocked,
				ErrorDescription: "Customer Locked",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp, nil
		}
	}
	resp := wa.PayResponse{
		Type:             wa.SyncBillPayResponse,
		TxnID:            request.TxnID,
		RefID:            refid,
		Result:           tigosdk.SucceededTxnResult,
		ErrorCode:        wa.ErrSuccessTxn,
		ErrorDescription: "Transaction Successful",
		Msisdn:           request.Msisdn,
		Flag:             tigosdk.YesFlag,
		Content:          request.SenderName,
	}
	return resp, nil
}

func (c checker) HandleSubscriberNameQuery(ctx context.Context, request wa.NameRequest) (wa.NameResponse, error) {

	user, found := c.checkUser(request.CustomerReferenceID)
	if !found {
		resp := wa.NameResponse{
			Type:      wa.SyncLookupResponse,
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
				Type:      wa.SyncLookupResponse,
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
				Type:      wa.SyncLookupResponse,
				Result:    "TF",
				ErrorCode: wa.ErrNameInvalidFormat,
				ErrorDesc: "Transaction Failed: Format not known",
				Msisdn:    request.Msisdn,
				Flag:      tigosdk.NoFlag,
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp, nil
		}

		resp := wa.NameResponse{
			Type:      wa.SyncLookupResponse,
			Result:    tigosdk.SucceededTxnResult,
			ErrorCode: wa.NoNamecheckErr,
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

func Server() *http.Server {

	app := makeApp()
	server := &http.Server{
		Addr:              ":8090",
		Handler:           app.makeHandler(),
		ReadTimeout:       60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		WriteTimeout:      60 * time.Second,
	}

	return server
}
