package services

import (
	"context"
	"testing"

	"github.com/eduardo/classicCarSearch/internal/models"
)

func TestMockSheetsService_GetAllParts(t *testing.T) {
	svc := NewMockSheetsService(nil)

	parts, err := svc.GetAllParts(context.Background())
	if err != nil {
		t.Fatalf("GetAllParts() error = %v", err)
	}

	if len(parts) == 0 {
		t.Error("GetAllParts() returned empty list")
	}

	for _, p := range parts {
		if p.Estado == "eliminado" {
			t.Errorf("GetAllParts() returned part with estado 'eliminado': %s", p.Name)
		}
	}
}

func TestMockSheetsService_GetFilteredParts(t *testing.T) {
	svc := NewMockSheetsService(nil)

	tests := []struct {
		name     string
		filters  models.FilterOptions
		expected int
	}{
		{
			name:     "no filters",
			filters:  models.FilterOptions{},
			expected: 21,
		},
		{
			name:     "filter by brand Ford",
			filters:  models.FilterOptions{Brand: "Ford"},
			expected: 4,
		},
		{
			name:     "filter by type Motor",
			filters:  models.FilterOptions{Type: "Motor"},
			expected: 8,
		},
		{
			name:     "combined filters",
			filters:  models.FilterOptions{Brand: "Ford", Type: "Enfriamiento"},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts, err := svc.GetFilteredParts(context.Background(), tt.filters)
			if err != nil {
				t.Fatalf("GetFilteredParts() error = %v", err)
			}
			if len(parts) != tt.expected {
				t.Errorf("GetFilteredParts() = %d parts, want %d", len(parts), tt.expected)
			}
		})
	}
}

func TestMockSheetsService_GetUniqueBrands(t *testing.T) {
	svc := NewMockSheetsService(nil)

	brands, err := svc.GetUniqueBrands(context.Background())
	if err != nil {
		t.Fatalf("GetUniqueBrands() error = %v", err)
	}

	if len(brands) == 0 {
		t.Error("GetUniqueBrands() returned empty list")
	}

	for _, b := range brands {
		if b == "" {
			t.Error("GetUniqueBrands() returned empty brand name")
		}
	}
}

func TestMockSheetsService_GetUsers(t *testing.T) {
	svc := NewMockSheetsService(nil)

	users, err := svc.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("GetUsers() error = %v", err)
	}

	if len(users) == 0 {
		t.Error("GetUsers() returned empty list")
	}

	found := false
	for _, u := range users {
		if u.Username == "admin" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetUsers() should contain 'admin' user")
	}
}

func TestMockSheetsService_CustomData(t *testing.T) {
	customData := &MockData{
		Parts: []models.Part{
			{ID: "1", Name: "Test Part", Brand: "Test", Type: "Test", Price: 100},
		},
		Users: []models.User{
			{Username: "testuser", Password: "testpass"},
		},
	}

	svc := NewMockSheetsService(customData)

	parts, _ := svc.GetAllParts(context.Background())
	if len(parts) != 1 {
		t.Errorf("GetAllParts() = %d parts, want 1", len(parts))
	}

	users, _ := svc.GetUsers(context.Background())
	if len(users) != 1 {
		t.Errorf("GetUsers() = %d users, want 1", len(users))
	}
}

func TestMockSheetsService_Close(t *testing.T) {
	svc := NewMockSheetsService(nil)
	if err := svc.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestNewDataProvider_MockMode(t *testing.T) {
	provider, err := NewDataProvider(true, "", "")
	if err != nil {
		t.Fatalf("NewDataProvider(true) error = %v", err)
	}

	_, ok := provider.(*MockSheetsService)
	if !ok {
		t.Error("NewDataProvider(true) should return MockSheetsService")
	}
}
