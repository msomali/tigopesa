# TIGO PUSH API

```code
TIGO_USERNAME=""
TIGO_PASSWORD=""
TIGO_BILLER_MSISDN=""
TIGO_BILLER_CODE=""
TIGO_GET_TOKEN_URL=""
TIGO_BILL_URL=""
TIGO_BASE_URL=""
DEBUG=true
```


### NOTE
- all TIGO_*_URL must be prefixed with / except for TIGO_BASE_URL.
- TIGO_BASE_URL format <ip_addr:port>.
- callback url is **/tigopesa/pushpay/callback**.

then run these commands

```bash
go build -o push && ./push
```

 initiate pushpay by sending request to localhost:8090/tigopesa/pushpay
```json
{
  "customer": 255xxxxxxxxx,
  "amount": 13000,
  "remarks":"tigo push pay"
}
```