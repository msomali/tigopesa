package push

const (
	ErrCodeCallbackSuccess = "BILLER-30-0000-S"
	ErrCodeCallbackFailed  = "BILLER-30-3030-E"
)

type (
	TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	PayRequest struct {
		CustomerMSISDN string `json:"CustomerMSISDN"`
		BillerMSISDN   string `json:"BillerMSISDN"`
		Amount         int    `json:"Amount"`
		Remarks        string `json:"Remarks,omitempty"`
		ReferenceID    string `json:"ReferenceID"`
	}

	PayResponse struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseStatus      bool   `json:"ResponseStatus"`
		ResponseDescription string `json:"ResponseDescription"`
		ReferenceID         string `json:"ReferenceID"`
		Message             string `json:"Message,omitempty"`
	}

	CallbackRequest struct {
		Status           bool   `json:"Status"`
		Description      string `json:"Description"`
		MFSTransactionID string `json:"MFSTransactionID,omitempty"`
		ReferenceID      string `json:"ReferenceID"`
		Amount           string `json:"Amount"`
	}

	CallbackResponse struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseDescription string `json:"ResponseDescription"`
		ResponseStatus      bool   `json:"ResponseStatus"`
		ReferenceID         string `json:"ReferenceID"`
	}

	RefundRequest struct {
		CustomerMSISDN      string `json:"CustomerMSISDN"`
		ChannelMSISDN       string `json:"ChannelMSISDN"`
		ChannelPIN          string `json:"ChannelPIN"`
		Amount              int    `json:"Amount"`
		MFSTransactionID    string `json:"MFSTransactionID"`
		ReferenceID         string `json:"ReferenceID"`
		PurchaseReferenceID string `json:"PurchaseReferenceID"`
	}

	RefundResponse struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseStatus      bool   `json:"ResponseStatus"`
		ResponseDescription string `json:"ResponseDescription"`
		DMReferenceID       string `json:"DMReferenceID"`
		ReferenceID         string `json:"ReferenceID"`
		MFSTransactionID    string `json:"MFSTransactionID,omitempty"`
		Message             string `json:"Message,omitempty"`
	}

	HealthCheckResponse struct {
		ReferenceID string `json:"ReferenceID"`
		Description string `json:"Description"`
	}
	HealthCheckRequest struct {
		ReferenceID string `json:"ReferenceID"`
	}
)
