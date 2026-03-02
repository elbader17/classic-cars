package services

import (
	"context"

	"github.com/eduardo/classicCarSearch/internal/models"
)

type DataProvider interface {
	GetAllParts(ctx context.Context) ([]models.Part, error)
	GetFilteredParts(ctx context.Context, filters models.FilterOptions) ([]models.Part, error)
	GetUniqueBrands(ctx context.Context) ([]string, error)
	GetUniqueCategories(ctx context.Context) ([]string, error)
	GetSubcategoriasByCategoria(ctx context.Context) (map[string][]string, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	Close() error
}

func NewDataProvider(mockMode bool, credentialsPath, spreadsheetID string) (DataProvider, error) {
	if mockMode {
		return NewMockSheetsService(nil), nil
	}

	return NewSheetsService(credentialsPath, spreadsheetID)
}
