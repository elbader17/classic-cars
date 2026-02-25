package services

import (
	"context"
	"strings"

	"github.com/eduardo/classicCarSearch/internal/models"
)

type MockData struct {
	Parts []models.Part
	Users []models.User
}

func DefaultMockData() *MockData {
	price1 := 450.00
	price2 := 320.00
	price3 := 150.00
	price4 := 1200.00
	price5 := 580.00
	price6 := 220.00
	price7 := 380.00
	price8 := 2500.00
	price9 := 340.00
	price10 := 890.00
	price11 := 420.00
	price12 := 180.00
	price13 := 95.00
	price14 := 480.00
	price15 := 150.00
	price16 := 650.00
	price17 := 45.00
	price18 := 220.00
	price19 := 320.00
	price20 := 1800.00
	price21 := 580.00
	price22 := 950.00
	return &MockData{
		Parts: []models.Part{
			{ID: "1", Name: "Carburador Ford V8", Brand: "Ford", Type: "Motor", Model: "Mustang", Year: "1967", Price: &price1, Description: "Carburador original en buen estado", Estado: "", Imagenes: "https://images.unsplash.com/photo-1486262715619-670810aa70f2?auto=format&fit=crop&q=80&w=1200,https://images.unsplash.com/photo-1590212151175-e58edd96185b?auto=format&fit=crop&q=80&w=1200"},
			{ID: "2", Name: "Radiador Chevrolet", Brand: "Chevrolet", Type: "Enfriamiento", Model: "Camaro", Year: "1969", Price: &price2, Description: "Radiador reconstruido", Estado: "", Imagenes: "https://images.unsplash.com/photo-1610647752706-3bb12232b3ab?auto=format&fit=crop&q=80&w=1200"},
			{ID: "3", Name: "Bomba de agua Ford", Brand: "Ford", Type: "Enfriamiento", Model: "Falcon", Year: "1970", Price: &price3, Description: "Nueva, en caja original", Estado: "", Imagenes: ""},
			{ID: "4", Name: "Diferencial Chrysler", Brand: "Chrysler", Type: "Transmision", Model: "Charger", Year: "1969", Price: &price4, Description: "Diferencial positraccion", Estado: "", Imagenes: ""},
			{ID: "5", Name: "Carburador Holley 4 bocas", Brand: "Holley", Type: "Motor", Model: "Universal", Year: "1970", Price: &price5, Description: "Carburador de alto rendimiento", Estado: "", Imagenes: ""},
			{ID: "6", Name: "Volante deportivo Momo", Brand: "Momo", Type: "Interior", Model: "Universal", Year: "1975", Price: &price6, Description: "Volante de cuero original", Estado: "", Imagenes: "https://images.unsplash.com/photo-1542282088-fe8426682b8f?auto=format&fit=crop&q=80&w=1200"},
			{ID: "7", Name: "Escape deportivo Ford", Brand: "Ford", Type: "Escape", Model: "Mustang", Year: "1968", Price: &price7, Description: "Sistema completo de escape", Estado: "eliminado", Imagenes: ""},
			{ID: "8", Name: "Caja de cambios Muncie", Brand: "Chevrolet", Type: "Transmision", Model: "Camaro", Year: "1969", Price: &price8, Description: "Transmision de 4 velocidades", Estado: "", Imagenes: ""},
			{ID: "9", Name: "Distribuidor electronico", Brand: "MSD", Type: "Motor", Model: "Universal", Year: "1972", Price: &price9, Description: "Sistema de encendido electronico", Estado: "", Imagenes: ""},
			{ID: "10", Name: "Tren delantero Mustang", Brand: "Ford", Type: "Suspension", Model: "Mustang", Year: "1966", Price: &price10, Description: "Brazos y rotulas completas", Estado: "", Imagenes: ""},
			{ID: "11", Name: "Carburador Weber doble", Brand: "Weber", Type: "Motor", Model: "Universal", Year: "1968", Price: &price11, Description: "Carburador italiano de alto flujo", Estado: "", Imagenes: ""},
			{ID: "12", Name: "Amortiguadores Koni", Brand: "Koni", Type: "Suspension", Model: "Universal", Year: "1970", Price: &price12, Description: "Par de amortiguadores traseros", Estado: "", Imagenes: ""},
			{ID: "13", Name: "Filtro de aire cromado", Brand: "Edelbrock", Type: "Motor", Model: "Universal", Year: "1969", Price: &price13, Description: "Filtro conico de alto flujo", Estado: "", Imagenes: ""},
			{ID: "14", Name: "Radiador aluminum Mustang", Brand: "Ford", Type: "Enfriamiento", Model: "Mustang", Year: "1965", Price: &price14, Description: "Radiador de aluminio racing", Estado: "", Imagenes: ""},
			{ID: "15", Name: "Parlantes Pioneer 6x9", Brand: "Pioneer", Type: "Interior", Model: "Universal", Year: "1975", Price: &price15, Description: "Par de parlantes coaxiales", Estado: "", Imagenes: ""},
			{ID: "16", Name: "Multiple de escape Hooker", Brand: "Hooker", Type: "Escape", Model: "Camaro", Year: "1970", Price: &price16, Description: "Headers de acero inoxidable", Estado: "", Imagenes: ""},
			{ID: "17", Name: "Bujias NGK iridium", Brand: "NGK", Type: "Motor", Model: "Universal", Year: "1970", Price: &price17, Description: "Set de 8 bujias", Estado: "", Imagenes: ""},
			{ID: "18", Name: "Faro delantero Chevrolet", Brand: "Chevrolet", Type: "Electrico", Model: "Camaro", Year: "1968", Price: &price18, Description: "Faro sellado original", Estado: "", Imagenes: ""},
			{ID: "19", Name: "Alternador Powermaster", Brand: "Powermaster", Type: "Electrico", Model: "Universal", Year: "1969", Price: &price19, Description: "Alternador 100 amperios", Estado: "", Imagenes: ""},
			{ID: "20", Name: "Asientos baquet Recaro", Brand: "Recaro", Type: "Interior", Model: "Universal", Year: "1975", Price: &price20, Description: "Par de asientos de competicion", Estado: "", Imagenes: ""},
			{ID: "21", Name: "Piston Forged JE", Brand: "JE Pistons", Type: "Motor", Model: "Ford V8", Year: "1968", Price: &price21, Description: "Set de 8 pistones forjados", Estado: "", Imagenes: ""},
			{ID: "22", Name: "Cigueñal Eagle", Brand: "Eagle", Type: "Motor", Model: "Ford V8", Year: "1969", Price: &price22, Description: "Cigueñal de acero forjado", Estado: "", Imagenes: ""},
		},
		Users: []models.User{
			{Username: "admin", Password: "admin123"},
			{Username: "usuario", Password: "test"},
		},
	}
}

type MockSheetsService struct {
	data *MockData
}

func NewMockSheetsService(data *MockData) *MockSheetsService {
	if data == nil {
		data = DefaultMockData()
	}
	return &MockSheetsService{data: data}
}

func (m *MockSheetsService) GetAllParts(ctx context.Context) ([]models.Part, error) {
	parts := filterActiveParts(m.data.Parts)
	for i := range parts {
		parsePartImages(&parts[i])
	}
	return parts, nil
}

func (m *MockSheetsService) GetFilteredParts(ctx context.Context, filters models.FilterOptions) ([]models.Part, error) {
	parts := filterActiveParts(m.data.Parts)
	parts = filterByBrand(parts, filters.Brand)
	return filterByType(parts, filters.Type), nil
}

func (m *MockSheetsService) GetUniqueBrands(ctx context.Context) ([]string, error) {
	parts := filterActiveParts(m.data.Parts)
	return uniqueValues(parts, func(p models.Part) string { return p.Brand }), nil
}

func (m *MockSheetsService) GetUniqueTypes(ctx context.Context) ([]string, error) {
	parts := filterActiveParts(m.data.Parts)
	return uniqueValues(parts, func(p models.Part) string { return p.Type }), nil
}

func (m *MockSheetsService) GetUsers(ctx context.Context) ([]models.User, error) {
	return m.data.Users, nil
}

func (m *MockSheetsService) Close() error {
	return nil
}

func filterByBrand(parts []models.Part, brand string) []models.Part {
	if brand == "" {
		return parts
	}
	filtered := make([]models.Part, 0, len(parts))
	for _, p := range parts {
		if p.Brand == brand {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func filterByType(parts []models.Part, partType string) []models.Part {
	if partType == "" {
		return parts
	}
	filtered := make([]models.Part, 0, len(parts))
	for _, p := range parts {
		if p.Type == partType {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func filterActiveParts(parts []models.Part) []models.Part {
	var result []models.Part
	for _, p := range parts {
		if p.Estado == "" || p.Estado != "eliminado" {
			result = append(result, p)
		}
	}
	return result
}

func uniqueValues[T any](items []T, fn func(T) string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		value := fn(item)
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func parsePartImages(part *models.Part) {
	if part.Imagenes == "" {
		part.ImagenesArr = nil
		return
	}
	part.ImagenesArr = splitAndTrim(part.Imagenes, ",")
}

func splitAndTrim(s, sep string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
