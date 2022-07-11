package entity

// Wallet is an entity to record wallet balance
type Wallet struct {
	WalletId string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}
