package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	sdk "github.com/techcraftlabs/tigopesa"
	"github.com/techcraftlabs/tigopesa/disburse"
	"github.com/techcraftlabs/tigopesa/examples"
	"github.com/techcraftlabs/tigopesa/internal"
	"github.com/techcraftlabs/tigopesa/push"
	"github.com/techcraftlabs/tigopesa/ussd"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	_ push.CallbackHandler  = (*Handler)(nil)
	_ ussd.NameQueryHandler = (*Handler)(nil)
	_ ussd.PaymentHandler   = (*Handler)(nil)
)

type (
	User struct {
		Name   string `json:"name"`
		RefID  string `json:"ref_id"`
		Status int    `json:"status"`
	}

	App struct {
		*sdk.Client
		Port              string
		ReadTimeout       time.Duration
		WriteTimeout      time.Duration
		ReadHeaderTimeout time.Duration
	}

	Handler struct {
		Users map[string]User
	}

	request struct {
		Msisdn  string  `json:"msisdn"`
		Amount  float64 `json:"amount"`
		Remarks string  `json:"remarks,omitempty"`
	}
)

func (h Handler) PaymentRequest(ctx context.Context, request ussd.PayRequest) (ussd.PayResponse, error) {
	user, found := h.checkUser(request.CustomerReferenceID)
	refid := fmt.Sprintf("%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	if !found {
		resp := ussd.PayResponse{

			Type:             ussd.SyncBillPayResponse,
			TxnID:            request.TxnID,
			RefID:            refid,
			Result:           sdk.FailedTxnResult,
			ErrorCode:        disburse.ErrInvalidCustomerRefNumber,
			ErrorDescription: "User Not Found",
			Msisdn:           request.Msisdn,
			Flag:             sdk.NoFlag,
			Content:          request.SenderName,
		}

		return resp, nil
	} else {
		if user.Status == 1 {
			resp := ussd.PayResponse{
				Type:             ussd.SyncBillPayResponse,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        ussd.ErrInvalidCustomerRefNumber,
				ErrorDescription: "Invalid Customer ref Number",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp, nil
		} else if user.Status == 2 {
			resp := ussd.PayResponse{
				Type:             sdk.SyncBillPayResponse,
				TxnID:            request.TxnID,
				RefID:            refid,
				Result:           "TF",
				ErrorCode:        ussd.ErrCustomerRefNumLocked,
				ErrorDescription: "Customer Locked",
				Msisdn:           request.Msisdn,
				Flag:             "N",
				Content:          request.SenderName,
			}
			return resp, nil
		}
	}
	resp := ussd.PayResponse{
		Type:             ussd.SyncBillPayResponse,
		TxnID:            request.TxnID,
		RefID:            refid,
		Result:           sdk.SucceededTxnResult,
		ErrorCode:        ussd.ErrSuccessTxn,
		ErrorDescription: "Transaction Successful",
		Msisdn:           request.Msisdn,
		Flag:             sdk.YesFlag,
		Content:          request.SenderName,
	}
	return resp, nil
}

func (h Handler) NameQuery(ctx context.Context, request ussd.NameRequest) (ussd.NameResponse, error) {
	user, found := h.checkUser(request.CustomerReferenceID)
	if !found {
		resp := ussd.NameResponse{
			Type:      ussd.SyncLookupResponse,
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

			resp := ussd.NameResponse{
				Type:      ussd.SyncLookupResponse,
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
			resp := ussd.NameResponse{
				Type:      ussd.SyncLookupResponse,
				Result:    "TF",
				ErrorCode: ussd.ErrNameInvalidFormat,
				ErrorDesc: "Transaction Failed: Format not known",
				Msisdn:    request.Msisdn,
				Flag:      sdk.NoFlag,
				Content:   fmt.Sprintf("%s", user.Name),
			}

			return resp, nil
		}

		resp := ussd.NameResponse{
			Type:      ussd.SyncLookupResponse,
			Result:    sdk.SucceededTxnResult,
			ErrorCode: ussd.NoNamecheckErr,
			ErrorDesc: "Transaction Successfully",
			Msisdn:    request.Msisdn,
			Flag:      sdk.YesFlag,
			Content:   fmt.Sprintf("%s", user.Name),
		}
		return resp, nil
	}
}

func (h Handler) Do(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {

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

func (h Handler) checkUser(ref string) (User, bool) {
	fmt.Printf("checking %s\n", ref)

	user, found := h.Users[ref]

	return user, found
}

func Server(app *App) *http.Server {

	server := &http.Server{
		Addr:              app.Port,
		Handler:           app.Handler(),
		ReadTimeout:       app.ReadTimeout,
		ReadHeaderTimeout: app.ReadHeaderTimeout,
		WriteTimeout:      app.WriteTimeout,
	}

	return server
}

func (app *App) DisburseHandler(writer http.ResponseWriter, r *http.Request) {
	var info request
	//	// Try to decode the r body into the struct. If there is an error,
	// respond to the pkg with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	ref := fmt.Sprintf("PCT%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	disburseRequest := disburse.Request{
		ReferenceID: ref,
		MSISDN:      info.Msisdn,
		Amount:      info.Amount,
	}
	resp, err := app.sendRequest(context.TODO(), internal.DISBURSE_REQUEST, disburseRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	//logger.Printf("disburse response %v\n",resp)

	writer.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(writer).Encode(resp)
}

func (app *App) PushPayHandler(writer http.ResponseWriter, r *http.Request) {
	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	billRequest := push.PayRequest{
		CustomerMSISDN: req.Msisdn,
		Amount:         int(req.Amount),
		Remarks:        req.Remarks,
		ReferenceID:    fmt.Sprintf("%s%d", app.Config.PushBillerCode, time.Now().Local().Unix()),
	}

	response, err := app.sendRequest(context.Background(), internal.PushPayRequest, billRequest)
	if err != nil {
		log.Printf("PushBillPay r failed error: %s", err.Error())
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(response)
}

func (app *App) Handler() http.Handler {
	ctx := context.TODO()
	router := mux.NewRouter()
	router.HandleFunc(app.Config.PayNamecheckURL,
		app.handleRequest(context.TODO(), internal.NameQueryRequest)).
		Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc(app.Config.PayRequestURL,
		app.handleRequest(ctx, internal.PaymentRequest)).
		Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc("/tigopesa/callback",
		app.handleRequest(context.TODO(), internal.CallbackRequest)).
		Methods(http.MethodPost)

	router.HandleFunc("/api/tigopesa/disburse", app.DisburseHandler).Methods(http.MethodPost)
	router.HandleFunc("/tigopesa/pushpay", app.PushPayHandler).Methods(http.MethodPost)

	return router
}

func main() {
	err := godotenv.Load("tigo.env")
	if err != nil {
		log.Printf("error %v\n", err)
		log.Fatal("Error loading .env file")
	}

	config := examples.LoadConfFromEnv()
	base := internal.NewBaseClient(internal.WithDebugMode(true))

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

	handler := Handler{
		Users: usersMap,
	}

	client := sdk.NewClient(config, base, handler, handler, handler)

	app := &App{
		Client:            client,
		Port:              ":8090",
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
	}
	err = Server(app).ListenAndServe()
	if err != nil {
		panic(err)
	}
}
