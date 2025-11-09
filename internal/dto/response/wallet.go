package response

// CheckWalletResponse represents the response for wallet existence check
type CheckWalletResponse struct {
	Exists    bool   `json:"exists"`
	AccountID string `json:"account_id,omitempty"`
}

// GetBalanceResponse represents the response for wallet balance
// Balance is in dirams (1 TJS = 100 dirams)
type GetBalanceResponse struct {
	AccountID string `json:"account_id"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
}
