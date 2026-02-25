package services

import (
	"testing"

	"github.com/eduardo/classicCarSearch/internal/models"
)

func TestFuzzySearch(t *testing.T) {
	parts := []models.Part{
		{ID: "1", Name: "Carburador Ford V8"},
		{ID: "2", Name: "Radiador Chevrolet"},
		{ID: "3", Name: "Bomba de agua"},
		{ID: "4", Name: "Carburador Holley"},
	}

	svc := NewSearchService()

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "exact match",
			query:         "Carburador",
			expectedCount: 2,
		},
		{
			name:          "partial match",
			query:         "carb",
			expectedCount: 2,
		},
		{
			name:          "typo tolerance",
			query:         "Carburadr",
			expectedCount: 2,
		},
		{
			name:          "no match",
			query:         "xyz123",
			expectedCount: 0,
		},
		{
			name:          "empty query returns all",
			query:         "",
			expectedCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := svc.FuzzySearch(tt.query, parts)
			if len(results) != tt.expectedCount {
				t.Errorf("FuzzySearch(%q) = %d results, want %d", tt.query, len(results), tt.expectedCount)
			}
		})
	}
}

func TestFuzzySearchWithFilters(t *testing.T) {
	parts := []models.Part{
		{ID: "1", Name: "Carburador Ford V8", Brand: "Ford", Type: "Motor", Price: 500},
		{ID: "2", Name: "Radiador Chevrolet", Brand: "Chevrolet", Type: "Enfriamiento", Price: 300},
		{ID: "3", Name: "Bomba de agua Ford", Brand: "Ford", Type: "Enfriamiento", Price: 150},
	}

	svc := NewSearchService()

	tests := []struct {
		name     string
		query    string
		filters  models.FilterOptions
		expected int
	}{
		{
			name:     "filter by brand",
			query:    "",
			filters:  models.FilterOptions{Brand: "Ford"},
			expected: 2,
		},
		{
			name:     "filter by type",
			query:    "",
			filters:  models.FilterOptions{Type: "Motor"},
			expected: 1,
		},
		{
			name:     "combined filters",
			query:    "",
			filters:  models.FilterOptions{Brand: "Ford", Type: "Enfriamiento"},
			expected: 1,
		},
		{
			name:     "search with filter",
			query:    "Bomba",
			filters:  models.FilterOptions{Type: "Enfriamiento"},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := svc.FuzzySearchWithFilters(tt.query, parts, tt.filters)
			if len(results) != tt.expected {
				t.Errorf("FuzzySearchWithFilters() = %d results, want %d", len(results), tt.expected)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		target   string
		distance int
		minScore int
		maxScore int
	}{
		{
			name:     "exact match",
			query:    "ford",
			target:   "ford",
			distance: 0,
			minScore: 100,
			maxScore: 100,
		},
		{
			name:     "partial match with bonus",
			query:    "for",
			target:   "ford",
			distance: 1,
			minScore: 80,
			maxScore: 100,
		},
		{
			name:     "contains bonus",
			query:    "ord",
			target:   "ford",
			distance: 1,
			minScore: 70,
			maxScore: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculateScore(tt.query, tt.target, tt.distance)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("calculateScore() = %d, want between %d and %d", score, tt.minScore, tt.maxScore)
			}
		})
	}
}
