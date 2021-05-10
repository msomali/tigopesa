# tigo
fully compliant tigo integration library in Go

## example
Follow these procedures to run the ussd example

go to `examples/ussd` folder create file `tigo.env` if not available and fill the proper values of each variable

```go
TIGO_USERNAME            = "TIGO_USERNAME"
TIGO_PASSWORD            = "TIGO_PASSWORD"
TIGO_PASSWORD_GRANT_TYPE = "TIGO_PASSWORD_GRANT_TYPE"
TIGO_ACCOUNT_NAME        = "TIGO_ACCOUNT_NAME"
TIGO_ACCOUNT_MSISDN      = "TIGO_ACCOUNT_MSISDN"
TIGO_BRAND_ID            = "TIGO_BRAND_ID"
TIGO_BILLER_CODE         = "TIGO_BILLER_CODE"
TIGO_GET_TOKEN_URL       = "TIGO_GET_TOKEN_URL"
TIGO_BILL_URL            = "TIGO_BILL_URL"
TIGO_A2W_URL             = "TIGO_A2W_URL"
TIGO_NAMECHECK_URL       = "TIGO_NAMECHECK_URL"
TIGO_W2A_URL             = "TIGO_W2A_URL"
TIGO_DISBURSEMENT_PIN    = "TIGO_DISBURSEMENT_PIN"
```
then run these commands

```bash
go build -o ussd && ./ussd

```