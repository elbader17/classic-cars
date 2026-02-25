package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eduardo/classicCarSearch/internal/models"
	"github.com/eduardo/classicCarSearch/internal/services"
)

// API Response types
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type PartsResponse struct {
	Results         []models.SearchResult `json:"results"`
	AvailableBrands []string              `json:"availableBrands"`
	AvailableTypes  []string              `json:"availableTypes"`
}

type APIHandler struct {
	provider services.DataProvider
	search   *services.SearchService
	auth     *services.AuthService
	session  *services.SessionService
}

func NewAPIHandler(provider services.DataProvider, search *services.SearchService, auth *services.AuthService, session *services.SessionService) *APIHandler {
	return &APIHandler{
		provider: provider,
		search:   search,
		auth:     auth,
		session:  session,
	}
}

func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSONError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if !h.auth.Authenticate(ctx, req.Username, req.Password) {
		writeJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := h.session.Create(req.Username)
	if err != nil {
		writeJSONError(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	h.session.SetCookie(w, token)

	writeJSONSuccess(w, LoginResponse{
		Token:    token,
		Username: req.Username,
	})
}

func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := h.session.GetTokenFromRequest(r)
	if token != "" {
		h.session.Delete(token)
	}

	h.session.ClearCookie(w)
	writeJSONSuccess(w, nil)
}

func (h *APIHandler) GetParts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	token := h.session.GetTokenFromRequest(r)
	if token == "" {
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session := h.session.Validate(token)
	if session == nil {
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	// Parse query parameters
	query := r.URL.Query().Get("q")
	brand := r.URL.Query().Get("brand")
	partType := r.URL.Query().Get("type")

	filters := models.FilterOptions{
		Brand: brand,
		Type:  partType,
	}

	parts, err := h.provider.GetAllParts(ctx)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := h.search.FuzzySearchWithFilters(query, parts, filters)

	// Get filter options
	brands, err := h.provider.GetUniqueBrands(ctx)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	types, err := h.provider.GetUniqueTypes(ctx)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(w, PartsResponse{
		Results:         results,
		AvailableBrands: brands,
		AvailableTypes:  types,
	})
}

func (h *APIHandler) GetFilters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check authentication
	token := h.session.GetTokenFromRequest(r)
	if token == "" {
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session := h.session.Validate(token)
	if session == nil {
		writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	brands, err := h.provider.GetUniqueBrands(ctx)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	types, err := h.provider.GetUniqueTypes(ctx)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(w, map[string][]string{
		"brands": brands,
		"types":  types,
	})
}

func writeJSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}
