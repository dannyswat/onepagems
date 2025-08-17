package types

import (
	"context"
	"time"
)

// SessionContextKey is the key used to store session in context
type SessionContextKey string

const SessionKey SessionContextKey = "session"

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}

// SessionContext creates a new context with the session
func SessionContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, SessionKey, session)
}

// SessionFromContext retrieves the session from context
func SessionFromContext(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(SessionKey).(*Session)
	return session, ok
}

// IsExpired checks if a session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt) || !s.IsActive
}

// Extend extends the session expiration time
func (s *Session) Extend(minutes int) {
	s.ExpiresAt = time.Now().Add(time.Duration(minutes) * time.Minute)
}
