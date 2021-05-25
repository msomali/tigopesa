package sdk

import (
	"context"
	"io"
	"net/http"
	"time"
)

type (
	Config struct {
		Username                     string
		Password                     string
		PasswordGrantType            string
		AccountName                  string
		AccountMSISDN                string
		BrandID                      string
		BillerCode                   string
		BillerMSISDN                 string
		ApiBaseURL                   string
		GetTokenRequestURL           string
		PushPayBillRequestURL        string
		PushPayReverseTransactionURL string
		PushPayHealthCheckURL        string
		AccountToWalletRequestURL    string
		AccountToWalletRequestPIN    string
		WalletToAccountRequestURL    string
		NameCheckRequestURL          string
	}

	Client struct {
		Config
		httpClient          *http.Client
		ctx                 context.Context
		timeout             time.Duration
		logger              io.Writer // for logging purposes
	}

	// ClientOption is a setter func to set Client details like
	// timeout, context, httpClient and logger
	ClientOption func(client *Client)

)
