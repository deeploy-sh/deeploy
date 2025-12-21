package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"sync"
	"time"
)

// CLISession represents a pending CLI authentication session
type CLISession struct {
	Token     string
	ExpiresAt time.Time
}

var (
	sessions   = make(map[string]*CLISession)
	sessionsMu sync.RWMutex
	sessionTTL = 5 * time.Minute
)

// CreateSession creates a new CLI auth session and returns its ID
func CreateSession() string {
	b := make([]byte, 16)
	rand.Read(b)
	sessionID := hex.EncodeToString(b)

	sessionsMu.Lock()
	sessions[sessionID] = &CLISession{
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	sessionsMu.Unlock()

	return sessionID
}

// SetSessionToken stores the auth token for a session (creates if not exists)
func SetSessionToken(sessionID, token string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	log.Printf("[CLI Auth] SetSessionToken: session=%s token=%s...", sessionID, token[:min(10, len(token))])

	session, exists := sessions[sessionID]
	if !exists {
		session = &CLISession{ExpiresAt: time.Now().Add(sessionTTL)}
		sessions[sessionID] = session
		log.Printf("[CLI Auth] Created new session: %s", sessionID)
	}

	session.Token = token
}

// GetSessionToken retrieves the token for a session
// Creates the session if it doesn't exist (TUI polls before login completes)
func GetSessionToken(sessionID string) (token string, ready bool) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	session, exists := sessions[sessionID]

	// Create session if it doesn't exist (TUI started polling)
	if !exists {
		sessions[sessionID] = &CLISession{ExpiresAt: time.Now().Add(sessionTTL)}
		log.Printf("[CLI Auth] GetSessionToken: session=%s CREATED (pending)", sessionID)
		return "", false
	}

	// Session expired
	if time.Now().After(session.ExpiresAt) {
		log.Printf("[CLI Auth] GetSessionToken: session=%s EXPIRED", sessionID)
		delete(sessions, sessionID)
		return "", false
	}

	// Token not set yet
	if session.Token == "" {
		log.Printf("[CLI Auth] GetSessionToken: session=%s PENDING", sessionID)
		return "", false
	}

	log.Printf("[CLI Auth] GetSessionToken: session=%s READY", sessionID)
	return session.Token, true
}

// DeleteSession removes a session (call after successful poll)
func DeleteSession(sessionID string) {
	sessionsMu.Lock()
	delete(sessions, sessionID)
	sessionsMu.Unlock()
}
