package app

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewSession(id uuid.UUID, username string, refreshToken, userAgent, clientIP string, IsBlocked bool, expiresAt time.Time) *Session {
	return &Session{
		ID:           id,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    IsBlocked,
		ExpiresAt:    expiresAt,
	}
}
