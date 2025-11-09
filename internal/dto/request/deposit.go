package request

// DepositRequest represents the request to deposit money into a wallet
// Amount is in dirams (1 TJS = 100 dirams)
type DepositRequest struct {
	AccountID string `json:"account_id" validate:"required,min=3,max=50"`
	Amount    int64  `json:"amount" validate:"required,gt=0"`
}
