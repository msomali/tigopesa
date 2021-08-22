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

package tigopesa

const (
	TXN_STATUS_200     = "200"
	TXN_STATUS_00026   = "00026"
	TXN_STATUS_00031   = "00031"
	TXN_STATUS_00042   = "00042"
	TXN_STATUS_317     = "317"
	TXN_STATUS_410     = "410"
	TXN_STATUS_2117    = "2117"
	TXN_STATUS_00317   = "00317"
	TXN_STATUS_00410   = "00410"
	TXN_STATUS_02117   = "02117"
	TXN_STATUS_0       = "0"
	TXN_STATUS_60014   = "60014"
	TXN_STATUS_60017   = "60017"
	TXN_STATUS_60018   = "60018"
	TXN_STATUS_60019   = "60019"
	TXN_STATUS_60021   = "60021"
	TXN_STATUS_60024   = "60024"
	TXN_STATUS_60028   = "60028"
	TXN_STATUS_60030   = "60030"
	TXN_STATUS_60074   = "60074"
	TXN_STATUS_100     = "100"
	YesFlag            = "Y"
	NoFlag             = "N"
	SucceededTxnResult = "TS"
	FailedTxnResult    = "TF"

	SyncLookupResponse  = "SYNC_LOOKUP_RESPONSE"
	SyncBillPayResponse = "SYNC_BILLPAY_RESPONSE"
	REQMFCI             = "REQMFCI"
)

var (
	txnStatusMap = map[string]string{
		"00026": "PIN expired. Please change your PIN.",
		"00031": "Requested amount more then allowed in the network",
		"00042": "Requested amount not in multiple of allowed value",
		"317":   "Unable to complete transaction as recipient A/c is barred. Error code 00317.",
		"410":   "Unable to complete transaction as amount is more than the maximum limit. Error code: 00410.",
		"2117":  "Unable to complete transaction as sender A/c is barred. Error code 02117",
		"00317": "Unable to complete transaction as recipient A/c is barred. Error code 00317.",
		"00410": "Unable to complete transaction as amount is more than the maximum limit. Error code: 00410.",
		"02117": "Unable to complete transaction as sender A/c is barred. Error code 02117",
		"200":   "Success",
		"0":     "Success",
		"60014": "Unable to complete transaction as maximum transaction value per day for payer reached. Error code 60014",
		"60017": "Unable to complete transaction as transaction amount is less than the minimum txn value for sender. Error code 60017. ",
		"60018": "Unable to complete transaction as amount is more than the maximum limit. Error code 60018. ",
		"60019": "Unable to complete transaction as account would go below minimum balance. Error code 60019. ",
		"60021": "Unable to complete transaction as maximum number of transactions per day for Payee was reached. Error code 60021. ",
		"60024": "Unable to complete transaction as maximum transaction value per day reached. Error code 60024. ",
		"60028": "Unable to complete transaction as transaction amount is more than the maximum txn value for recipient. Error code 60028.",
		"60030": "Unable to complete transaction as the Payee account would go above maximum balance. Error code: 60030.",
		"60074": "Payee Role Type Transfer Profile not defined",
		"100":   "This is generic error, which is returned if problem happen during transaction\nprocessing. Partner should put transaction amount in HOLD state to avoid risk of\nrollback while amount was disbursed. This is the same case for any kind of\ntimeout as well.\n",
	}
)

func TxnStatusDesc(errCode string) (string, bool) {
	desc, exists := txnStatusMap[errCode]
	return desc, exists
}
