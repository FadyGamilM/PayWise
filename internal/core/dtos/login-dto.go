package dtos

import (
	"time"

	"github.com/google/uuid"
)

type LoginReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRes struct {
	SessionID              uuid.UUID   `json:"session_id"`
	AccessToken            string      `json:"access_token"`
	AccessTokenExpiration  time.Time   `json:"access_token_expires_at"`
	RefreshToken           string      `json:"refresh_token"`
	RefreshTokenExpiration time.Time   `json:"refresh_token_expires_at"`
	User                   *UserResDto `json:"user"`
}
