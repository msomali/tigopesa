# Integration credentials (TIGO)
A compilation of all credentials and environmental variables
that are used for the TigoPesa integration

## USSD

### Tigo Pesa USSD (i.e. *150​*01# ) C2B or App
```
Account Name: TECH CRAFT
Account MSISDN: 25565298765
Biller Number/Company number: 298765
TIGO_W2A_URL: http://96.126.126.124:8090/api/tigopesa/c2b/transactions
TIGO_NAMECHECK_URL: http://96.126.126.124:8090/api/tigopesa/c2b/users/names

```
### B2C
```
Account Name: TECH CRAFT MFI
Account MSISDN: 25564000892
Brand ID: 2341
TIGO_A2W_URL: http://accessgwtest.tigo.co.tz:8080/TECHCRAFT2Tigo
TIGO_DISBURSEMENT_PIN: 7645

```

## PUSH 

```java
TIGO_USERNAME="TECHCRAFT"
TIGO_PASSWORD="idBnPwe"
TIGO_BILLER_MSISDN="25565298765"
TIGO_BILLER_CODE="TCR"
TIGO_GET_TOKEN_URL="http://accessgwtest.tigo.co.tz:8080/TechCraft2DM-GetToken"
TIGO_PUSH_URL: "http://accessgwtest.tigo.co.tz:8080/TechCraft2DM-PushBillpay"
TIGO_PUSH_CALLBACK_URL:  "http://96.126.126.124:8090/tigopesa/pushpay/callback"
```