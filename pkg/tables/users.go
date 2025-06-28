package tables

import "time"

type User struct {
	UserId          int64      `json:"userId"`
	UserCode        string     `json:"userCode"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	FullName        string     `json:"fullName"`
	Username        *string    `json:"username,omitempty"`
	PenName         *string    `json:"penName,omitempty"`
	AvatarURL       string     `json:"avatarUrl"`
	LoginWith       string     `json:"loginWith"`
	IsEmailVerified bool       `json:"isEmailVerified"`
	Phone           *string    `json:"phone,omitempty"`
	Instagram       *string    `json:"instagram,omitempty"`
	BankId          *int64     `json:"bankId,omitempty"`
	AccountNumber   *string    `json:"accountNumber,omitempty"`
	FlgAuthor       string     `json:"flgAuthor"`
	CreateDatetime  time.Time  `json:"createDatetime"`
	UpdateDatetime  *time.Time `json:"updateDatetime,omitempty"`
}
