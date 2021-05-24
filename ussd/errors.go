package ussd

var TxnStatus = map[string]string{
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

const (
	NAMECHECK_SUCCESS        = "error000"
	NAMECHECK_NOT_REGISTERED = "error010"
	NAMECHECK_SUSPENDED      = "error020"
	NAMECHECK_INVALID_FORMAT = "error030"
)

const (
	W2A_ERR_SUCCESSFUL_TRANSACTION      = "error000"
	W2A_ERR_SERVICE_NOT_AVAILABLE       = "error001"
	W2A_ERR_INVALID_CUSTOMER_REF_NUMBER = "error010"
	W2A_ERR_CUSTOMER_REF_NUM_LOCKED     = "error011"
	W2A_ERR_INVALID_AMOUNT              = "error012"
	W2A_ERR_AMOUNT_INSUFFICIENT         = "error013"
	W2A_ERR_AMOUNT_TOO_HIGH             = "error014"
	W2A_ERR_AMOUNT_TOO_LOW              = "error015"
	W2A_ERR_INVALID_PAYMENT             = "error016"
	W2A_ERR_GENERAL_ERROR               = "error100"
	W2A_ERR_RETRY_CONDITION_NO_RESPONSE = "error111"
)
