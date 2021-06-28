package conf

import (
	"errors"
	"github.com/techcraftt/tigosdk/aw"
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/wa"
)

var (
	_ Validator = (*Config)(nil)

	ErrConfigNil = errors.New("config can not be nil")
)

type (
	Config struct {
		PayAccountName   string
		PayAccountMSISDN string
		PayBillerNumber  string
		PayRequestURL    string
		PayNamecheckURL  string

		DisburseAccountName   string
		DisburseAccountMSISDN string
		DisburseBrandID       string
		DisbursePIN           string
		DisburseRequestURL    string

		PushUsername              string
		PushPassword              string
		PushPasswordGrantType     string
		PushApiBaseURL            string
		PushGetTokenURL           string
		PushBillerMSISDN          string
		PushBillerCode            string
		PushPayURL                string
		PushReverseTransactionURL string
		PushHealthCheckURL        string
	}

	//Validator validates the configurations and return nil if all
	//is good or return an error
	Validator interface {
		Validate() error
	}
)

func (conf *Config) Validate() error {
	panic("implement me")
}

func (conf *Config) Split() (pushConf *push.Config, pay *wa.Config, disburse *aw.Config) {
	pushConf = &push.Config{
		Username:              conf.PushUsername,
		Password:              conf.PushPassword,
		PasswordGrantType:     conf.PushPasswordGrantType,
		ApiBaseURL:            conf.PushApiBaseURL,
		GetTokenURL:           conf.PushGetTokenURL,
		BillerMSISDN:          conf.PushBillerMSISDN,
		BillerCode:            conf.PushBillerCode,
		PushPayURL:            conf.PushPayURL,
		ReverseTransactionURL: conf.PushReverseTransactionURL,
		HealthCheckURL:        conf.PushHealthCheckURL,
	}

	pay = &wa.Config{
		AccountName:   conf.PayAccountName,
		AccountMSISDN: conf.PayAccountMSISDN,
		BillerNumber:  conf.PayBillerNumber,
		RequestURL:    conf.PayRequestURL,
		NamecheckURL:  conf.PayNamecheckURL,
	}

	disburse = &aw.Config{
		AccountName:   conf.DisburseAccountName,
		AccountMSISDN: conf.DisburseAccountMSISDN,
		BrandID:       conf.DisburseBrandID,
		PIN:           conf.DisbursePIN,
		RequestURL:    conf.DisburseRequestURL,
	}

	return
}

// Merge combine configurations of different clients. It is usefully when they have been loaded from
// different sources before being used:
// returns error ErrConfigNil if any of the 3 config is nil
func Merge(pushConf *push.Config, waConf *wa.Config, awConf *aw.Config) (*Config, error) {

	if pushConf == nil || waConf == nil || awConf == nil {
		return nil, ErrConfigNil
	}

	merged := &Config{
		PayAccountName:            waConf.AccountName,
		PayAccountMSISDN:          waConf.AccountMSISDN,
		PayBillerNumber:           waConf.BillerNumber,
		PayRequestURL:             waConf.RequestURL,
		PayNamecheckURL:           waConf.NamecheckURL,
		DisburseAccountName:       awConf.AccountName,
		DisburseAccountMSISDN:     awConf.AccountMSISDN,
		DisburseBrandID:           awConf.BrandID,
		DisbursePIN:               awConf.PIN,
		DisburseRequestURL:        awConf.RequestURL,
		PushUsername:              pushConf.Username,
		PushPassword:              pushConf.Password,
		PushPasswordGrantType:     pushConf.PasswordGrantType,
		PushApiBaseURL:            pushConf.ApiBaseURL,
		PushGetTokenURL:           pushConf.GetTokenURL,
		PushBillerMSISDN:          pushConf.BillerMSISDN,
		PushBillerCode:            pushConf.BillerCode,
		PushPayURL:                pushConf.PushPayURL,
		PushReverseTransactionURL: pushConf.ReverseTransactionURL,
		PushHealthCheckURL:        pushConf.HealthCheckURL,
	}
	return merged, nil
}
