package managers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"onepagems/internal/types"
	"time"
)

// AuthManager handles authentication and session management
type AuthManager struct {
	sessions map[string]*types.Session
	config   *types.Config
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(config *types.Config) *AuthManager {
	return &AuthManager{
		sessions: make(map[string]*types.Session),
		config:   config,
	}
}

// Login authenticates a user and creates a session
func (am *AuthManager) Login(username, password string) (*types.Session, error) {
	// Hash the provided password
	hashedPassword := am.hashPassword(password)

	// Check against configured credentials
	if username != am.config.AdminUsername || hashedPassword != am.config.AdminPassword {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create new session
	sessionID, err := am.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	session := &types.Session{
		ID:        sessionID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour sessions
		IsActive:  true,
	}

	am.sessions[sessionID] = session
	return session, nil
}

// Logout invalidates a session
func (am *AuthManager) Logout(sessionID string) error {
	if session, exists := am.sessions[sessionID]; exists {
		session.IsActive = false
		delete(am.sessions, sessionID)
		return nil
	}
	return fmt.Errorf("session not found")
}

// ValidateSession checks if a session is valid and active
func (am *AuthManager) ValidateSession(sessionID string) (*types.Session, error) {
	session, exists := am.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	if !session.IsActive {
		return nil, fmt.Errorf("session is inactive")
	}

	if time.Now().After(session.ExpiresAt) {
		session.IsActive = false
		delete(am.sessions, sessionID)
		return nil, fmt.Errorf("session has expired")
	}

	// Extend session expiry on successful validation
	session.ExpiresAt = time.Now().Add(24 * time.Hour)

	return session, nil
}

// GetSessionFromRequest extracts session ID from HTTP request
func (am *AuthManager) GetSessionFromRequest(r *http.Request) (*types.Session, error) {
	// Try to get session ID from cookie first
	cookie, err := r.Cookie("session_id")
	if err == nil {
		return am.ValidateSession(cookie.Value)
	}

	// Fall back to Authorization header
	sessionID := r.Header.Get("Authorization")
	if sessionID == "" {
		return nil, fmt.Errorf("no session ID provided")
	}

	// Remove "Bearer " prefix if present
	if len(sessionID) > 7 && sessionID[:7] == "Bearer " {
		sessionID = sessionID[7:]
	}

	return am.ValidateSession(sessionID)
}

// RequireAuth is a middleware that requires authentication
func (am *AuthManager) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := am.GetSessionFromRequest(r)
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Add session to request context
		r = r.WithContext(types.SessionContext(r.Context(), session))
		next(w, r)
	}
}

// CreateSessionCookie creates an HTTP cookie for the session
func (am *AuthManager) CreateSessionCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours in seconds
	}
}

// ClearSessionCookie creates a cookie that clears the session
func (am *AuthManager) ClearSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // Delete cookie
	}
}

// CleanupExpiredSessions removes expired sessions from memory
func (am *AuthManager) CleanupExpiredSessions() {
	now := time.Now()
	for sessionID, session := range am.sessions {
		if now.After(session.ExpiresAt) || !session.IsActive {
			delete(am.sessions, sessionID)
		}
	}
}

// GetActiveSessions returns the count of active sessions
func (am *AuthManager) GetActiveSessions() int {
	am.CleanupExpiredSessions()
	return len(am.sessions)
}

// ListSessions returns all active sessions (for admin purposes)
func (am *AuthManager) ListSessions() []*types.Session {
	am.CleanupExpiredSessions()
	sessions := make([]*types.Session, 0, len(am.sessions))
	for _, session := range am.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// generateSessionID generates a cryptographically secure session ID
func (am *AuthManager) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashPassword creates a SHA-256 hash of the password
func (am *AuthManager) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// ChangePassword changes the admin password (requires current password)
func (am *AuthManager) ChangePassword(currentPassword, newPassword string) error {
	currentHashed := am.hashPassword(currentPassword)
	if currentHashed != am.config.AdminPassword {
		return fmt.Errorf("current password is incorrect")
	}

	if len(newPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters long")
	}

	// Update the config with new hashed password
	am.config.AdminPassword = am.hashPassword(newPassword)
	return nil
}
