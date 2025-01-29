package model

type TransactionRequest struct {
	Amount  float64 `json:"amount"`
	Remarks string  `json:"remarks"`
	Type    string
	UserAccess
}

type TransactionResponse struct {
	TopUpID       string  `json:"top_up_id,omitempty"`
	PaymentID     string  `json:"payment_id,omitempty"`
	Amount        float64 `json:"amount"`
	Remarks       string  `json:"remarks,omitempty"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	CreatedDate   string  `json:"created_date"`
}
