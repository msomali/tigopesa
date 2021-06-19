# tigosdk
tigosdk is open source fully compliant tigo pesa client written in golang


## contents
1. [usage](#usage)
2. [example](#example)
2. [projects](#projects)
3. [links](#links)
4. [contributors](#contributors)
5. [sponsors](#sponsers)

## usage
```bash

go get https://github.com/techcraftt/tigosdk

```

```go
import (
   "github.com/techcraftt/tigosdk"
)
```


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
**TO ENABLE DEBUG MODE CALL WithDebugMode(true) before With HttpClient**


## projects
The List of projects using this library

1. [PAYCRAFT]() - Full-Fledged Payment as a Service

## links

## contributors

1. [Bethuel Mmbaga]()
2. [Frances Ruganyumisa]()
3. [Pius Alfred]()

## sponsors

[Techcraft Technologies Co. LTD]()
