# tigo
fully compliant tigo integration library in Go

## example
Follow these procedures to run the example

add correct details to files available at `cmd/main.go`
```go

conf := tigosdk.Configs{
		Username:          "",
		Password:          "",
		PasswordGrantType: "",
		AccountName:       "",
		AccountMSISDN:     "",
		BrandID:           "",
		BillerCode:        "",
		GetTokenURL:       "",
		BillURL:           "",
		A2WReqURL:         "",
	}

```

then

```bash

cd cmd && go run main.go

```