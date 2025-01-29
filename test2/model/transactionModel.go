package model

type TransactionRequest struct {
	TargetUser string  `json:"target_user"`
	Amount     float64 `json:"amount"`
	Remarks    string  `json:"remarks"`
	Type       string
	UserAccess
}

type TransactionResponse struct {
	TopUpID       string  `json:"top_up_id,omitempty"`
	PaymentID     string  `json:"payment_id,omitempty"`
	TransferID    string  `json:"transfer_id,omitempty"`
	Amount        float64 `json:"amount"`
	Remarks       string  `json:"remarks,omitempty"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	CreatedDate   string  `json:"created_date"`
}
