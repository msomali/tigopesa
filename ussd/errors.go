package ussd

type NameCheckErrorCode string

const (
	SUCCESS        NameCheckErrorCode = "error000"
	NOT_REGISTERED NameCheckErrorCode = "error010"
	SUSPENDED      NameCheckErrorCode = "error020"
	INVALID_FORMAT NameCheckErrorCode = "error030"
)

type W2AErrorCode string

const (
	SUCCESSFUL_TRANSACTION      W2AErrorCode = "error000"
	SERVICE_NOT_AVAILABLE       W2AErrorCode = "error001"
	INVALID_CUSTOMER_REF_NUMBER W2AErrorCode = "error010"
	CUSTOMER_REF_NUM_LOCKED     W2AErrorCode = "error011"
	INVALID_AMOUNT              W2AErrorCode = "error012"
	AMOUNT_INSUFFICIENT         W2AErrorCode = "error013"
	AMOUNT_TOO_HIGH             W2AErrorCode = "error014"
	AMOUNT_TOO_LOW              W2AErrorCode = "error015"
	INVALID_PAYMENT             W2AErrorCode = "error016"
	GENERAL_ERROR               W2AErrorCode = "error100"
	RETRY_CONDITION_NO_RESPONSE W2AErrorCode = "error111"
)
