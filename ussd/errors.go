package ussd

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
