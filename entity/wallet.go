package entity

type Wallet struct {
	WalletId string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}
