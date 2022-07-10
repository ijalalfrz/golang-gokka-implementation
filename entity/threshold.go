package entity

// Threshold is an entitiy to record deposit and check threshold
type Threshold struct {
	WalletId                 string  `json:"wallet_id"`
	Deposit                  float64 `json:"deposit"`
	TotalDepositWithinWindow float64 `json:"total_deposit_within_window"`
	StartWindowTime          int64   `json:"start_window_time"`
	CreatedTime              int64   `json:"created_time"`
	AboveThreshold           bool    `json:"above_threshold"`
}
