package webmodel

type DepositWalletPayload struct {
	WalletId string  `json:"wallet_id" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
}
