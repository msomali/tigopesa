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

## NOTE
- all URL must be prefixed with / except for BASE URL.
- BASE URL format <ip_addr>:<port>

then run these commands

```bash
go build -o push && ./push
```

