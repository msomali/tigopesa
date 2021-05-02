package tigo

//Config acontains details of TigoPesa integration
//These are configurations supplied during the integration stage
type Config struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	PasswordGrantType string `json:"password_grant_type"`
	AccountName       string `json:"account_name"`
	AccountMSISDN     string `json:"account_msisdn"`
	BrandID           string `json:"brand_id"`
	BillerCode        string `json:"biller_code"`
	GetTokenURL       string `json:"get_token_url"`
	BillURL           string `json:"biller_payment"`
}
