package webmodel

type DepositWalletPayload struct {
	WalletId string  `json:"wallet_id" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
}

type DetailWalletResponse struct {
	WalletId       string  `json:"wallet_id"`
	Balance        float64 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}
