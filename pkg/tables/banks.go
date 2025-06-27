package tables

import "time"

type Bank struct {
	BankID            int64      `json:"bank_id"`
	BankName          string     `json:"bank_name"`
	BankCode          *string    `json:"bank_code,omitempty"`
	Remark            *string    `json:"remark,omitempty"`
	ActiveDatetime    time.Time  `json:"active_datetime"`
	NonActiveDatetime *time.Time `json:"non_active_datetime,omitempty"`
	CreateDatetime    time.Time  `json:"create_datetime"`
	UpdateDatetime    *time.Time `json:"update_datetime,omitempty"`
}
