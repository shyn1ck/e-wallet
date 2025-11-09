package response

// DepositResponse represents the response for a deposit operation
// Amount and NewBalance are in dirams (1 TJS = 100 dirams)
type DepositResponse struct {
	Success       bool   `json:"success"`
	AccountID     string `json:"account_id"`
	Amount        int64  `json:"amount"`
	NewBalance    int64  `json:"new_balance"`
	Currency      string `json:"currency"`
	TransactionID int64  `json:"transaction_id"`
}
