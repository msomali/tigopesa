package push

type (
	BillPayResponse struct {
		ResponseCode        string `json:"ResponseCode"`
		ResponseStatus      bool   `json:"ResponseStatus"`
		ResponseDescription string `json:"ResponseDescription"`
		ReferenceID         string `json:"ReferenceID"`
		Message             string `json:"Message,omitempty"`
	}

	RefundPaymentResponse struct {
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
)
