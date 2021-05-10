package ussd

type Config struct {
	ServerBaseURL string
	ServerPort string
	NameCheckReqEndpoint string
	WalletToAccountReqEndpoint string
	DisbursementReqEndpoint string
	DisbursementPIN string

}
