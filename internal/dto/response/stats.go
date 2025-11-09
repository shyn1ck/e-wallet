package response

// MonthlyStatsResponse represents the response for monthly statistics
// TotalAmount is in dirams (1 TJS = 100 dirams)
type MonthlyStatsResponse struct {
	AccountID   string `json:"account_id"`
	Month       string `json:"month"`
	TotalCount  int64  `json:"total_count"`
	TotalAmount int64  `json:"total_amount"`
	Currency    string `json:"currency"`
}
