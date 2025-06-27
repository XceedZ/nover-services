package tables

import "time"

type User struct {
	UserID          int64     `json:"user_id"`
	UserCode        string    `json:"user_code"`
	Email           string    `json:"email"`
	Password        string    `json:"password"` 
	FullName        string    `json:"full_name"`
	Username        *string   `json:"username,omitempty"`
	AvatarURL       string    `json:"avatar_url"`
	LoginWith       string    `json:"login_with"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreateDatetime  time.Time `json:"create_datetime"`
	UpdateDatetime  *time.Time`json:"update_datetime,omitempty"`
}