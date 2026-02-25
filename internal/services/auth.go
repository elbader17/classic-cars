package services

import (
	"context"
	"crypto/subtle"
	"sync"

	"github.com/eduardo/classicCarSearch/internal/models"
)

type AuthService struct {
	usersCache []models.User
	cacheMu    sync.RWMutex
	provider   DataProvider
}

func NewAuthService(provider DataProvider) *AuthService {
	return &AuthService{
		provider: provider,
	}
}

func (a *AuthService) Authenticate(ctx context.Context, username, password string) bool {
	users, err := a.GetUsers(ctx)
	if err != nil {
		return false
	}

	for _, user := range users {
		if subtle.ConstantTimeCompare([]byte(user.Username), []byte(username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) == 1 {
			return true
		}
	}

	return false
}

func (a *AuthService) GetUsers(ctx context.Context) ([]models.User, error) {
	a.cacheMu.RLock()
	if a.usersCache != nil {
		a.cacheMu.RUnlock()
		return a.usersCache, nil
	}
	a.cacheMu.RUnlock()

	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()

	users, err := a.provider.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	// Fallback: add default user if no users found in sheet
	if len(users) == 0 {
		users = []models.User{
			{Username: "admin", Password: "admin123"},
		}
	}

	a.usersCache = users
	return users, nil
}

func (a *AuthService) ClearCache() {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()
	a.usersCache = nil
}
