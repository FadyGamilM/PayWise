package models

import "time"

type Session struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	ClientIP     string    `json:"client_ip"`
	UserAgent    string    `json:"user_agent"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpireAt     time.Time `json:"expire_at"`
}
