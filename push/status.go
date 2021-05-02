package push

var statusText = map[string]string{
	// successful status
	"BILLER-18-0000-S":            "Valid request",
	"RefundTransaction-20-0000-S": "Reverse Transaction Success",

	// failure status
	"DebitMandate-10-1000-E":      "Internal Service Error",
	"DebitMandate-10-2001-E":      "User name or password empty",
	"DebitMandate-10-2002-E":      "CustomerMSISDN not passed",
	"DebitMandate-10-2003-E":      "Password not passed",
	"DebitMandate-10-2004-V":      "Invalid MSISDN",
	"DebitMandate-10-2005-E":      "Not Registered in DM",
	"DebitMandate-10-5000-E":      "<Backend error description>",
	"DebitMandate-10-3000-E":      "Username/Password provided is wrong",
	"DebitMandate-10-2038-V":      "BillerMSISDN not passed",
	"BILLER-18-3040-E":            "BillerMSISDN is wrong",
	"BILLER-18-3018-E":            "Failed in Sum Amount",
	"BILLER-18-3019-E":            "Failed in Min and Max Amount",
	"BILLER-18-3020-E":            "Failed in Frequency",
	"BILLER-18-3021-E":            "Biller Username or Password is wrong",
	"BILLER-18-3022-E":            "Biller not active",
	"DebitMandate-10-2020-V":      " Not Valid Reference ID",
	"RefundTransaction-20-3039-E": "Reverse Transaction Failed",
}
