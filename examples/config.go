package examples

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/techcraftlabs/tigopesa/pkg/conf"
	"log"
	"os"
)

const (
	envPayAccountName        = "TIGO_PAY_ACCOUNT_NAME"
	envPayAccountMSISDN      = "TIGO_PAY_ACCOUNT_MSISDN"
	envPayBillerNumber       = "TIGO_PAY_BILLER_NUMBER"
	envPayRequestURL         = "TIGO_PAY_REQUEST_URL"
	envPayNamecheckURL       = "TIGO_PAY_NAMECHECK_URL"
	envDisburseAccountName   = "TIGO_DISBURSE_ACCOUNT_NAME"
	envDisburseAccountMSISDN = "TIGO_DISBURSE_ACCOUNT_MSISDN"
	envDisburseBrandID       = "TIGO_DISBURSE_BRAND_ID"
	envDisbursePin           = "TIGO_DISBURSE_PIN"
	envDisburseURL           = "TIGO_DISBURSE_URL"
	envPushUsername          = "TIGO_PUSH_USERNAME"
	envPushPassword          = "TIGO_PUSH_PASSWORD"
	envPushBaseURL           = "TIGO_PUSH_BASE_URL"
	envPushBillerMSISDN      = "TIGO_PUSH_BILLER_MSISDN"
	envPushBillerCode        = "TIGO_PUSH_BILLER_CODE"
	envPushGetTokenURL       = "TIGO_PUSH_TOKEN_URL"
	envPushPayURL            = "TIGO_PUSH_PAY_URL"
)

func LoadConfFromEnv() *conf.Config {
	err := godotenv.Load("tigo.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	payAccountName := os.Getenv(envPayAccountName)
	payAccountMSISDN := os.Getenv(envPayAccountMSISDN)
	payBillerNumber := os.Getenv(envPayBillerNumber)
	payRequestURL := os.Getenv(envPayRequestURL)
	namecheckURL := os.Getenv(envPayNamecheckURL)

	disburseName := os.Getenv(envDisburseAccountName)
	disburseMSISDN := os.Getenv(envDisburseAccountMSISDN)
	disburseBrandID := os.Getenv(envDisburseBrandID)
	disbursePIN := os.Getenv(envDisbursePin)
	disburseURL := os.Getenv(envDisburseURL)

	pushName := os.Getenv(envPushUsername)
	pushPassword := os.Getenv(envPushPassword)
	pushBaseURL := os.Getenv(envPushBaseURL)
	pushTokenURL := os.Getenv(envPushGetTokenURL)
	pushPayURL := os.Getenv(envPushPayURL)
	pushMSISDN := os.Getenv(envPushBillerMSISDN)
	pushBillerCode := os.Getenv(envPushBillerCode)

	fmt.Printf("%s\n", payAccountName)
	fmt.Printf("%s\n", payAccountMSISDN)
	fmt.Printf("%s\n", payBillerNumber)
	fmt.Printf("%s\n", payRequestURL)
	fmt.Printf("%s\n", namecheckURL)
	fmt.Printf("%s\n", disburseName)
	fmt.Printf("%s\n", disburseMSISDN)
	fmt.Printf("%s\n", disburseBrandID)
	fmt.Printf("%s\n", disbursePIN)
	fmt.Printf("%s\n", disburseURL)
	fmt.Printf("%s\n", pushName)
	fmt.Printf("%s\n", pushPassword)
	fmt.Printf("%s\n", pushBaseURL)
	fmt.Printf("%s\n", pushTokenURL)
	fmt.Printf("%s\n", pushMSISDN)
	fmt.Printf("%s\n", pushBillerCode)
	fmt.Printf("%s\n", pushPayURL)

	return &conf.Config{
		PayAccountName:            payAccountName,
		PayAccountMSISDN:          payAccountMSISDN,
		PayBillerNumber:           payBillerNumber,
		PayRequestURL:             payRequestURL,
		PayNamecheckURL:           namecheckURL,
		DisburseAccountName:       disburseName,
		DisburseAccountMSISDN:     disburseMSISDN,
		DisburseBrandID:           disburseBrandID,
		DisbursePIN:               disbursePIN,
		DisburseRequestURL:        disburseURL,
		PushUsername:              pushName,
		PushPassword:              pushPassword,
		PushPasswordGrantType:     "password",
		PushApiBaseURL:            pushBaseURL,
		PushGetTokenURL:           pushTokenURL,
		PushBillerMSISDN:          pushMSISDN,
		PushBillerCode:            pushBillerCode,
		PushPayURL:                pushPayURL,
		PushReverseTransactionURL: "",
		PushHealthCheckURL:        "",
	}
}
