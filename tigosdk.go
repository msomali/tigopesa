package tigosdk

const (
	SYNC_LOOKUP_RESPONSE  = "SYNC_LOOKUP_RESPONSE"
	SYNC_BILLPAY_RESPONSE = "SYNC_BILLPAY_RESPONSE"
)

//Configs contains details of TigoPesa integration
//These are configurations supplied during the integration stage
type Configs struct {
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	PasswordGrantType         string `json:"grant_type"`
	AccountName               string `json:"account_name"`
	AccountMSISDN             string `json:"account_msisdn"`
	BrandID                   string `json:"brand_id"`
	BillerCode                string `json:"biller_code"`
	GetTokenRequestURL        string `json:"token_url"`
	PushPayBillRequestURL     string `json:"bill_url"`
	AccountToWalletRequestURL string `json:"a2w_url"`
	WalletToAccountRequestURL string `json:"w2a_url"`
	NameCheckRequestURL       string `json:"name_url"`
}

