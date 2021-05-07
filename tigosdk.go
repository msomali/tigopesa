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
	PasswordGrantType         string `json:"password_grant_type"`
	AccountName               string `json:"account_name"`
	AccountMSISDN             string `json:"account_msisdn"`
	BrandID                   string `json:"brand_id"`
	BillerCode                string `json:"biller_code"`
	GetTokenRequestURL        string `json:"get_token_url"`
	PushPayBillRequestURL     string `json:"biller_payment"`
	AccountToWalletRequestURL string `json:"a2w_url"`
	WalletToAccountRequestURL string `json:"w_2_aurl"`
	NameCheckRequestURL       string `json:"name_check_url"`
}

