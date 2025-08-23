package tables

import "time"

// Wallet merepresentasikan record dalam tabel wallets.
type Wallet struct {
	WalletID       int64      `json:"-" db:"wallet_id"` // Biasanya tidak perlu ditampilkan
	UserID         int64      `json:"userId" db:"user_id"`
	PaidCoins      int64      `json:"paidCoins" db:"paid_coins"`
	BonusCoins     int64      `json:"bonusCoins" db:"bonus_coins"`
	UpdateDatetime *time.Time `json:"updateDatetime,omitempty" db:"update_datetime"`
}