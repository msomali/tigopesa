package push

type GetTokenResponse struct {
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type BillPayResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseStatus      bool   `json:"ResponseStatus"`
	ResponseDescription string `json:"ResponseDescription"`
	ReferenceID         string `json:"ReferenceID"`
	Message             string `json:"Message,omitempty"`
}

type RefundPaymentResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseStatus      bool   `json:"ResponseStatus"`
	ResponseDescription string `json:"ResponseDescription"`
	DMReferenceID       string `json:"DMReferenceID"`
	ReferenceID         string `json:"ReferenceID"`
	MFSTransactionID    string `json:"MFSTransactionID,omitempty"`
	Message             string `json:"Message,omitempty"`
}

type HealthCheckResponse struct {
	ReferenceID string `json:"ReferenceID"`
	Description string `json:"Description"`
}
