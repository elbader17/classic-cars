package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/eduardo/classicCarSearch/internal/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	client        *http.Client
	spreadsheetID string
	service       *sheets.Service
}

func NewSheetsService(credentialsPath, spreadsheetID string) (*SheetsService, error) {
	client, err := getClient(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to create sheets client: %v", err)
	}

	svc, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to create sheets service: %v", err)
	}

	return &SheetsService{
		client:        client,
		spreadsheetID: spreadsheetID,
		service:       svc,
	}, nil
}

func getClient(credentialsPath string) (*http.Client, error) {
	// Check for GOOGLE_CREDENTIALS env var (base64 encoded JSON)
	if credsJSON := os.Getenv("GOOGLE_CREDENTIALS"); credsJSON != "" {
		conf, err := google.JWTConfigFromJSON([]byte(credsJSON), sheets.SpreadsheetsReadonlyScope)
		if err != nil {
			return nil, fmt.Errorf("unable to parse credentials from env: %v", err)
		}
		return conf.Client(context.Background()), nil
	}

	// Fallback: read from file
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}

	return conf.Client(context.Background()), nil
}

func (s *SheetsService) GetAllParts(ctx context.Context) ([]models.Part, error) {
	resp, err := s.service.Spreadsheets.Values.Get(s.spreadsheetID, "Repuestos!A:Z").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) < 2 {
		return nil, fmt.Errorf("no data found in sheet")
	}

	headers := resp.Values[0]
	parts := make([]models.Part, 0, len(resp.Values)-1)

	for _, row := range resp.Values[1:] {
		part := parseRow(headers, row)
		if part.Estado != "eliminado" {
			parsePartImages(&part)
			parts = append(parts, part)
		}
	}

	return parts, nil
}

func (s *SheetsService) GetFilteredParts(ctx context.Context, filters models.FilterOptions) ([]models.Part, error) {
	parts, err := s.GetAllParts(ctx)
	if err != nil {
		return nil, err
	}

	if filters.Brand != "" {
		var filtered []models.Part
		for _, p := range parts {
			if p.Brand == filters.Brand {
				filtered = append(filtered, p)
			}
		}
		parts = filtered
	}

	if filters.Type != "" {
		var filtered []models.Part
		for _, p := range parts {
			if p.Type == filters.Type {
				filtered = append(filtered, p)
			}
		}
		parts = filtered
	}

	return parts, nil
}

func (s *SheetsService) GetUniqueBrands(ctx context.Context) ([]string, error) {
	parts, err := s.GetAllParts(ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var brands []string
	for _, p := range parts {
		if p.Brand != "" && !seen[p.Brand] {
			seen[p.Brand] = true
			brands = append(brands, p.Brand)
		}
	}
	return brands, nil
}

func (s *SheetsService) GetUniqueTypes(ctx context.Context) ([]string, error) {
	parts, err := s.GetAllParts(ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var types []string
	for _, p := range parts {
		if p.Type != "" && !seen[p.Type] {
			seen[p.Type] = true
			types = append(types, p.Type)
		}
	}
	return types, nil
}

func (s *SheetsService) GetUsers(ctx context.Context) ([]models.User, error) {
	resp, err := s.service.Spreadsheets.Values.Get(s.spreadsheetID, "Usuarios!A:B").Context(ctx).Do()
	if err != nil {
		return nil, nil
	}

	if len(resp.Values) < 2 {
		return nil, nil
	}

	users := make([]models.User, 0, len(resp.Values)-1)
	for _, row := range resp.Values[1:] {
		if len(row) >= 2 {
			users = append(users, models.User{
				Username: fmt.Sprintf("%v", row[0]),
				Password: fmt.Sprintf("%v", row[1]),
			})
		}
	}

	return users, nil
}

func (s *SheetsService) Close() error {
	return nil
}

func parseRow(headers []interface{}, row []interface{}) models.Part {
	part := models.Part{}
	for i, header := range headers {
		if i >= len(row) {
			break
		}
		value := fmt.Sprintf("%v", row[i])
		switch strings.ToLower(header.(string)) {
		case "id":
			part.ID = value
		case "nombre":
			part.Name = value
		case "marca":
			part.Brand = value
		case "tipo":
			part.Type = value
		case "modelo":
			part.Model = value
		case "año", "ano":
			part.Year = value
		case "precio":
			fmt.Sscanf(value, "%f", &part.Price)
		case "descripcion":
			part.Description = value
		case "imagenes":
			part.Imagenes = value
		case "estado":
			part.Estado = value
		}
	}
	return part
}

var _ = strings.ToLower
