package services

import (
	"testing"
	"time"
)

func TestSessionService_Create(t *testing.T) {
	svc := NewSessionService()

	token, err := svc.Create("testuser")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if token == "" {
		t.Error("Create() returned empty token")
	}
}

func TestSessionService_Validate(t *testing.T) {
	svc := NewSessionService()

	token, _ := svc.Create("testuser")

	session := svc.Validate(token)
	if session == nil {
		t.Error("Validate() returned nil for valid token")
	}
	if session.Username != "testuser" {
		t.Errorf("Validate() username = %q, want %q", session.Username, "testuser")
	}
}

func TestSessionService_Validate_InvalidToken(t *testing.T) {
	svc := NewSessionService()

	session := svc.Validate("invalid-token")
	if session != nil {
		t.Error("Validate() should return nil for invalid token")
	}
}

func TestSessionService_Delete(t *testing.T) {
	svc := NewSessionService()

	token, _ := svc.Create("testuser")
	svc.Delete(token)

	session := svc.Validate(token)
	if session != nil {
		t.Error("Validate() should return nil after Delete()")
	}
}

func TestSessionService_Validate_ExpiredSession(t *testing.T) {
	svc := &SessionService{
		sessions: make(map[string]*SessionData),
		duration: 24 * time.Hour,
	}

	token := "test-token"
	svc.sessions[token] = &SessionData{
		Username:  "testuser",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	session := svc.Validate(token)
	if session != nil {
		t.Error("Validate() should return nil for expired session")
	}
}

func TestSessionService_Cleanup(t *testing.T) {
	svc := &SessionService{
		sessions: make(map[string]*SessionData),
		duration: 24 * time.Hour,
	}

	svc.sessions["active"] = &SessionData{
		Username:  "activeuser",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	svc.sessions["expired"] = &SessionData{
		Username:  "expireduser",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	svc.Cleanup()

	if svc.Validate("active") == nil {
		t.Error("active session should still exist")
	}
	if svc.Validate("expired") != nil {
		t.Error("expired session should be removed")
	}
}
