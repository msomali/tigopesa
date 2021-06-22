package tigosdk

const (
	TXN_STATUS_200   = "200"
	TXN_STATUS_00026 = "00026"
	TXN_STATUS_00031 = "00031"
	TXN_STATUS_00042 = "00042"
	TXN_STATUS_317   = "317"
	TXN_STATUS_410   = "410"
	TXN_STATUS_2117  = "2117"
	TXN_STATUS_00317 = "00317"
	TXN_STATUS_00410 = "00410"
	TXN_STATUS_02117 = "02117"
	TXN_STATUS_0     = "0"
	TXN_STATUS_60014 = "60014"
	TXN_STATUS_60017 = "60017"
	TXN_STATUS_60018 = "60018"
	TXN_STATUS_60019 = "60019"
	TXN_STATUS_60021 = "60021"
	TXN_STATUS_60024 = "60024"
	TXN_STATUS_60028 = "60028"
	TXN_STATUS_60030 = "60030"
	TXN_STATUS_60074 = "60074"
	TXN_STATUS_100   = "100"

	ErrNameNotRegistered = "error010"
	ErrNameInvalidFormat = "error030"
	ErrNameUserSuspended = "error030"
	NoNamecheckErr       = "error000"

	YesFlag            = "Y"
	NoFlag             = "N"
	SucceededTxnResult = "TS"
	FailedTxnResult    = "TF"

	SyncLookupResponse  = "SYNC_LOOKUP_RESPONSE"
	SyncBillPayResponse = "SYNC_BILLPAY_RESPONSE"
	REQMFCI             = "REQMFCI"

	// Payment error codes

	ErrSuccessTxn               = "error000"
	ErrServiceNotAvailable      = "error001"
	ErrInvalidCustomerRefNumber = "error010"
	ErrCustomerRefNumLocked     = "error011"
	ErrInvalidAmount            = "error012"
	ErrAmountInsufficient       = "error013"
	ErrAmountTooHigh            = "error014"
	ErrAmountTooLow             = "error015"
	ErrInvalidPayment           = "error016"
	ErrGeneralError             = "error100"
	ErrRetryConditionNoResponse = "error111"
)

var (
	txnStatusMap = map[string]string{
		"00026": "PIN expired. Please change your PIN.",
		"00031": "Requested amount more then allowed in the network",
		"00042": "Requested amount not in multiple of allowed value",
		"317":   "Unable to complete transaction as recipient A/c is barred. Error code 00317.",
		"410":   "Unable to complete transaction as amount is more than the maximum limit. Error code: 00410.",
		"2117":  "Unable to complete transaction as sender A/c is barred. Error code 02117",
		"00317": "Unable to complete transaction as recipient A/c is barred. Error code 00317.",
		"00410": "Unable to complete transaction as amount is more than the maximum limit. Error code: 00410.",
		"02117": "Unable to complete transaction as sender A/c is barred. Error code 02117",
		"200":   "Success",
		"0":     "Success",
		"60014": "Unable to complete transaction as maximum transaction value per day for payer reached. Error code 60014",
		"60017": "Unable to complete transaction as transaction amount is less than the minimum txn value for sender. Error code 60017. ",
		"60018": "Unable to complete transaction as amount is more than the maximum limit. Error code 60018. ",
		"60019": "Unable to complete transaction as account would go below minimum balance. Error code 60019. ",
		"60021": "Unable to complete transaction as maximum number of transactions per day for Payee was reached. Error code 60021. ",
		"60024": "Unable to complete transaction as maximum transaction value per day reached. Error code 60024. ",
		"60028": "Unable to complete transaction as transaction amount is more than the maximum txn value for recipient. Error code 60028.",
		"60030": "Unable to complete transaction as the Payee account would go above maximum balance. Error code: 60030.",
		"60074": "Payee Role Type Transfer Profile not defined",
		"100":   "This is generic error, which is returned if problem happen during transaction\nprocessing. Partner should put transaction amount in HOLD state to avoid risk of\nrollback while amount was disbursed. This is the same case for any kind of\ntimeout as well.\n",
	}
)

type (
//SubscriberNameRequest struct {
//	XMLName             xml.Name `xml:"COMMAND"`
//	Text                string   `xml:",chardata"`
//	Type                string   `xml:"TYPE"`
//	Msisdn              string   `xml:"MSISDN"`
//	CompanyName         string   `xml:"COMPANYNAME"`
//	CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
//}
//
//SubscriberNameResponse struct {
//	XMLName   xml.Name `xml:"COMMAND"`
//	Text      string   `xml:",chardata"`
//	Type      string   `xml:"TYPE"`
//	Result    string   `xml:"RESULT"`
//	ErrorCode string   `xml:"ERRORCODE"`
//	ErrorDesc string   `xml:"ERRORDESC"`
//	Msisdn    string   `xml:"MSISDN"`
//	Flag      string   `xml:"FLAG"`
//	Content   string   `xml:"CONTENT"`
//}
//
//AccountToWalletRequest struct {
//	XMLName     xml.Name `xml:"COMMAND"`
//	Text        string   `xml:",chardata"`
//	Type        string   `xml:"TYPE"`
//	ReferenceID string   `xml:"REFERENCEID"`
//	Msisdn      string   `xml:"MSISDN"`
//	PIN         string   `xml:"PIN"`
//	Msisdn1     string   `xml:"MSISDN1"`
//	Amount      float64  `xml:"AMOUNT"`
//	SenderName  string   `xml:"SENDERNAME"`
//	Language1   string   `xml:"LANGUAGE1"`
//	BrandID     string   `xml:"BRAND_ID"`
//}
//
//AccountToWalletResponse struct {
//	XMLName     xml.Name `xml:"COMMAND" json:"-"`
//	Text        string   `xml:",chardata" json:"-"`
//	Type        string   `xml:"TYPE" json:"type"`
//	ReferenceID string   `xml:"REFERENCEID" json:"reference_id"`
//	TxnID       string   `xml:"TXNID" json:"txnid"`
//	TxnStatus   string   `xml:"TXNSTATUS" json:"txn_status"`
//	Message     string   `xml:"MESSAGE" json:"message"`
//}
//
//WalletToAccountRequest struct {
//	XMLName             xml.Name `xml:"COMMAND"`
//	Text                string   `xml:",chardata"`
//	TYPE                string   `xml:"TYPE"`
//	TxnID               string   `xml:"TXNID"`
//	Msisdn              string   `xml:"MSISDN"`
//	Amount              float64  `xml:"AMOUNT"`
//	CompanyName         string   `xml:"COMPANYNAME"`
//	CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
//	SenderName          string   `xml:"SENDERNAME"`
//}
//
//WalletToAccountResponse struct {
//	XMLName          xml.Name `xml:"COMMAND"`
//	Text             string   `xml:",chardata"`
//	Type             string   `xml:"TYPE"`
//	TxnID            string   `xml:"TXNID"`
//	RefID            string   `xml:"REFID"`
//	Result           string   `xml:"RESULT"`
//	ErrorCode        string   `xml:"ERRORCODE"`
//	ErrorDescription string   `xml:"ERRORDESCRIPTION"`
//	Msisdn           string   `xml:"MSISDN"`
//	Flag             string   `xml:"FLAG"`
//	Content          string   `xml:"CONTENT"`
//}
)

func TxnStatusDesc(errCode string) (string, bool) {
	desc, exists := txnStatusMap[errCode]
	return desc, exists
}
