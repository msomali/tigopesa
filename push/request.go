package push

type BillPayRequest struct {
	CustomerMSISDN int64  `json:"CustomerMSISDN"`
	BillerMSISDN   int64  `json:"BillerMSISDN"`
	Amount         int    `json:"Amount"`
	Remarks        string `json:"Remarks,omitempty"`
	ReferenceID    string `json:"ReferenceID"`
}

type BillPayCallbackRequest struct {
	Status           bool   `json:"Status"`
	Description      string `json:"Description"`
	MFSTransactionID string `json:"MFSTransactionID,omitempty"`
	ReferenceID      string `json:"ReferenceID"`
	Amount           string    `json:"Amount"`
}

type RefundPaymentRequest struct {
	CustomerMSISDN      int64  `json:"CustomerMSISDN"`
	ChannelMSISDN       int64  `json:"ChannelMSISDN"`
	ChannelPIN          string `json:"ChannelPIN"`
	Amount              int    `json:"Amount"`
	MFSTransactionID    string `json:"MFSTransactionID"`
	ReferenceID         string `json:"ReferenceID"`
	PurchaseReferenceID string `json:"PurchaseReferenceID"`
}
type HealthCheckRequest struct {
	ReferenceID string `json:"ReferenceID"`
}
