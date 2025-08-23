package tables

import "time"

type CoinTransaction struct {
	TransactionID   int64      `json:"transactionId" db:"transaction_id"`
	UserID          int64      `json:"-" db:"user_id"`
	TransactionType string     `json:"transactionType" db:"transaction_type"`
	CoinType        string     `json:"coinType" db:"coin_type"`
	Amount          int        `json:"amount" db:"amount"`
	Description     *string    `json:"description,omitempty" db:"description"`
	ExpiryDate      *time.Time `json:"expiryDate,omitempty" db:"expiry_date"`
	CreateDatetime  time.Time  `json:"createDatetime" db:"create_datetime"`
}

type PaginatedTransactionResponse struct {
    Pagination PaginationInfo      `json:"pagination"`
    Items      []CoinTransaction `json:"transactions"`
}