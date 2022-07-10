package webmodel

// DepositWalletPayload is model for deposit wallet http request payload
type DepositWalletPayload struct {
	WalletId string  `json:"wallet_id" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
}

// DetailWalletResponse is response for get detail wallet
type DetailWalletResponse struct {
	WalletId       string  `json:"wallet_id"`
	Balance        float64 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}
