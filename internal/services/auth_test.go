package services

import (
	"context"
	"testing"

	"github.com/eduardo/classicCarSearch/internal/models"
)

func TestAuthService_Authenticate(t *testing.T) {
	users := []models.User{
		{Username: "admin", Password: "secret123"},
		{Username: "user", Password: "password"},
	}

	tests := []struct {
		name     string
		username string
		password string
		users    []models.User
		expected bool
	}{
		{
			name:     "valid credentials",
			username: "admin",
			password: "secret123",
			users:    users,
			expected: true,
		},
		{
			name:     "invalid password",
			username: "admin",
			password: "wrong",
			users:    users,
			expected: false,
		},
		{
			name:     "invalid username",
			username: "nobody",
			password: "secret123",
			users:    users,
			expected: false,
		},
		{
			name:     "empty credentials",
			username: "",
			password: "",
			users:    users,
			expected: false,
		},
		{
			name:     "no users",
			username: "admin",
			password: "secret123",
			users:    []models.User{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &AuthService{usersCache: tt.users}
			result := auth.Authenticate(context.Background(), tt.username, tt.password)
			if result != tt.expected {
				t.Errorf("Authenticate(%q, %q) = %v, want %v", tt.username, tt.password, result, tt.expected)
			}
		})
	}
}

func TestAuthService_ClearCache(t *testing.T) {
	auth := &AuthService{
		usersCache: []models.User{{Username: "test", Password: "test"}},
	}

	auth.ClearCache()

	if auth.usersCache != nil {
		t.Error("ClearCache() did not clear the cache")
	}
}
