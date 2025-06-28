package tables

import "time"

type Bank struct {
	BankId            int64      `json:"bankId"`
	BankName          string     `json:"bankName"`
	BankCode          *string    `json:"bankCode,omitempty"`
	Remark            *string    `json:"remark,omitempty"`
	ActiveDatetime    time.Time  `json:"activeDatetime"`
	NonActiveDatetime *time.Time `json:"nonActiveDatetime,omitempty"`
	CreateDatetime    time.Time  `json:"createDatetime"`
	UpdateDatetime    *time.Time `json:"updateDatetime,omitempty"`
}
