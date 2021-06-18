# integration

During integration with tigo there are details that will be provided essential in developing the system.


## contents
- [push pay integration](#pushpay)
- [b2c integration](#b2c)
- [c2b integration](#c2b)


## pushpay
This is to enable the business to bill their customer from the web by sending push pay request.

Integration information given to business by tigo during integration
- Username
- Password
- Biller MSISDN
- Biller Code
- Get Token URL - use this to ask token from TigoPesa gateway
- Push URL - post to this URL to initiate push pay request


Information required by Tigo
- Callback URL - Tigo post results of push pay requests to this. 

## b2c
This is to enablebusiness to pay cutsomers via web. Another common name is `Disbursement` or `Account To Wallet`

Details offered by Tigo during integrations are:
- Account Name
- Account MSISDN
- Account To Wallet URL - push disbursement requests here
- Disbursement PIN - Your Password


## c2b
This is to enable customers paying you through USSD i.e `*150*01#` or via application. The other name of these type of transactions is `Wallet To Account`.

Information provided by Tigo are 
- Account Name
- Account MSISDN
- Biller Number or Company Number

Information Required by Tigo
- Namecheck URL. It used by Tigo to ask if the paying customer and other details are correct
- Tigo Wallet To Account URL


NOTE: Some credentials like Account Name and Account MSISDN are offered only once and can be used in every case per a single integration.
