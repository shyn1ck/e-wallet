package request

// CheckWalletRequest represents the request to check if a wallet exists
type CheckWalletRequest struct {
	AccountID string `json:"account_id" validate:"required,min=3,max=50"`
}

// GetBalanceRequest represents the request to get wallet balance
type GetBalanceRequest struct {
	AccountID string `json:"account_id" validate:"required,min=3,max=50"`
}

// GetMonthlyStatsRequest represents the request to get monthly statistics
type GetMonthlyStatsRequest struct {
	AccountID string `json:"account_id" validate:"required,min=3,max=50"`
}
