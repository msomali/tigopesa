// Package ussd contains everything you will need to craft a tigopesa ussd pkg
//import (
//	"context"
//	"log"
//	"net/http"
//	"time"
//)
//
//func doc(){
//	// But Before creating the pkg for ussd you must first create the config object
//	// which has all the details obtained during network integration stage with tigo
//	conf := Config{
//		Username:                  "",
//		Password:                  "",
//		AccountName:               "",
//		AccountMSISDN:             "",
//		BrandID:                   "",
//		BillerCode:                "",
//		AccountToWalletRequestURL: "",
//		AccountToWalletRequestPIN: "",
//		WalletToAccountRequestURL: "",
//		NameCheckRequestURL:       "",
//	}
//
//	// Then one should create a QuerySubscriberFunc that has all the logic needed
//	// to identify the validity of the SubscriberNameRequest from TigoPesa
//
//	var namesHandleFunc QuerySubscriberFunc
//
//	namesHandleFunc = func(ctx context.Context, request SubscriberNameRequest) (resp SubscriberNameResponse, err error) {
//		// Check the request and return appropriate response
//		// For example you may check the name submitted against your database to authenticate
//		// the user or even perform authorization
//
//		return resp, err
//	}
//
//	// Then create a WalletToAccountFunc a function responsible to check if the collection from
//	// user account should proceed or not. The request is first received by WalletToAccountHandler
//	// before it is passed to this WalletToAccountFunc
//	// it determines from WalletToAccountRequest if the amount to be collected is right
//	// or if the Reference ID used is correct then after those operation the function returns
//	// a response to WalletToAccountHandler where it is properly marshalled and returned to TigoPesa
//	var collectionHandler WalletToAccountFunc
//
//	collectionHandler = func(ctx context.Context, request WalletToAccountRequest) (r WalletToAccountResponse, e error) {
//		// check against db, auth system etc etc
//		return r,e
//	}
//
//	// You can create pkg with default context.Context, timeout, io.Writer and http.Client
//	pkg := NewClient(conf,collectionHandler,namesHandleFunc)
//
//	// Or you can rather roll out your own
//
//	var writer = log.Writer()
//	var timeout = time.Second *30
//	var httpClient = &http.Client{
//		Transport:     nil,
//		CheckRedirect: nil,
//		Jar:           nil,
//		Timeout:       0,
//	}
//
//	var ctx, _ = context.WithTimeout(context.Background(), time.Minute)
//
//	// now you can wire in your own values
//	client2 := NewClient(conf,collectionHandler,namesHandleFunc,
//		WithContext(ctx),WithTimeout(timeout),WithLogger(writer),WithHTTPClient(httpClient))
//
//
//	var mux http.ServeMux
//
//	mux.HandleFunc(pkg.NameCheckRequestURL,pkg.SubscriberNameHandler)
//	mux.HandleFunc(pkg.WalletToAccountRequestURL,pkg.WalletToAccountHandler)
//	mux.HandleFunc(pkg.AccountToWalletRequestURL, disburseHandler)
//}
//
//func disburseHandler(writer http.ResponseWriter, request *http.Request) {
//
//}
package ussd
