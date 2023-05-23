package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type SessionRepository interface {
	Create(ctx context.Context, arg *Session) error
	Get(ctx context.Context, id uuid.UUID) (*Session, error)
}

func NewSession(id uuid.UUID, userID int64, refreshToken, userAgent, clientIP string, IsBlocked bool, expiresAt time.Time) *Session {
	return &Session{
		ID:           id,
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    IsBlocked,
		ExpiresAt:    expiresAt,
	}
}
