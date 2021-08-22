/*
 * MIT License
 *
 * Copyright (c) 2021 TechCraft Technologies Co. Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package config

import (
	"errors"
	"github.com/techcraftlabs/tigopesa/disburse"
	"github.com/techcraftlabs/tigopesa/push"
	"github.com/techcraftlabs/tigopesa/ussd"
)

var (
	ErrConfigNil = errors.New("config can not be nil")
)

type (
	Overall struct {
		PayAccountName            string
		PayAccountMSISDN          string
		PayBillerNumber           string
		PayRequestURL             string
		PayNamecheckURL           string
		DisburseAccountName       string
		DisburseAccountMSISDN     string
		DisburseBrandID           string
		DisbursePIN               string
		DisburseRequestURL        string
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
)

func (conf *Overall) Split() (pushConf *push.Config, pay *ussd.Config, dConf *disburse.Config) {
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

	pay = &ussd.Config{
		AccountName:   conf.PayAccountName,
		AccountMSISDN: conf.PayAccountMSISDN,
		BillerNumber:  conf.PayBillerNumber,
		RequestURL:    conf.PayRequestURL,
		NamecheckURL:  conf.PayNamecheckURL,
	}

	dConf = &disburse.Config{
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
func Merge(pushConf *push.Config, waConf *ussd.Config, awConf *disburse.Config) (*Overall, error) {

	if pushConf == nil || waConf == nil || awConf == nil {
		return nil, ErrConfigNil
	}

	merged := &Overall{
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
