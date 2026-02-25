package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/eduardo/classicCarSearch/internal/handlers"
	"github.com/eduardo/classicCarSearch/internal/services"
)

func main() {
	mockMode := getEnvBool("MOCK_MODE", false)
	credentialsPath := getEnv("CREDENTIALS_PATH", "")
	spreadsheetID := getEnv("SPREADSHEET_ID", "")
	port := getEnv("PORT", "8080")

	if !mockMode && spreadsheetID == "" {
		log.Fatal("SPREADSHEET_ID environment variable is required (or set MOCK_MODE=true)")
	}

	provider, err := services.NewDataProvider(mockMode, credentialsPath, spreadsheetID)
	if err != nil {
		log.Fatalf("Failed to create data provider: %v", err)
	}
	defer provider.Close()

	searchSvc := services.NewSearchService()
	authSvc := services.NewAuthService(provider)
	sessionSvc := services.NewSessionService()

	apiHandler := handlers.NewAPIHandler(provider, searchSvc, authSvc, sessionSvc)

	mux := http.NewServeMux()

	// API endpoints (JSON)
	mux.HandleFunc("/api/login", apiHandler.Login)
	mux.HandleFunc("/api/logout", apiHandler.Logout)
	mux.HandleFunc("/api/parts", apiHandler.GetParts)
	mux.HandleFunc("/api/filters", apiHandler.GetFilters)
	mux.HandleFunc("/api/image", apiHandler.ProxyImage)

	// Static files (including index.html as frontend)
	fs := http.FileServer(http.Dir("."))
	mux.Handle("/", fs)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		mode := "production"
		if mockMode {
			mode = "mock"
		}
		log.Printf("Server starting on port %s (mode: %s)", port, mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return boolVal
}
