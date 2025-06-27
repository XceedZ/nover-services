package tables

import "time"

type User struct {
	UserID          int64      `json:"user_id"`
	UserCode        string     `json:"user_code"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	FullName        string     `json:"full_name"`
	Username        *string    `json:"username,omitempty"`
	PenName         *string    `json:"pen_name,omitempty"`
	AvatarURL       string     `json:"avatar_url"`
	LoginWith       string     `json:"login_with"`
	IsEmailVerified bool       `json:"is_email_verified"`
	Phone           *string    `json:"phone,omitempty"`
	Instagram       *string    `json:"instagram,omitempty"`
	BankID          *int64     `json:"bank_id,omitempty"`
	AccountNumber   *string    `json:"account_number,omitempty"`
	FlgAuthor       string     `json:"flg_author"`
	CreateDatetime  time.Time  `json:"create_datetime"`
	UpdateDatetime  *time.Time `json:"update_datetime,omitempty"`
}
