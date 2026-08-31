package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	routes "padi-template/app/Routes"
	_ "padi-template/database/migrations"
	"github.com/wibiesana/padi_go_core/config"
	"github.com/wibiesana/padi_go_core/database"
	"github.com/wibiesana/padi_go_core/migrator"
	"github.com/wibiesana/padi_go_core/router"
)

func setupTestApp(t *testing.T) *router.Router {
	cfg := &config.Config{
		AppName:         "Padi Test API",
		AppEnv:          "testing",
		AppPort:         "8080",
		AppDebug:        true,
		AppKey:          "testing-key",
		JWTSecret:       "test-secret-32-chars-long-minimum-jwt",
		JWTExpiration:   1,
		DBConnection:    "sqlite",
		DBDatabase:      ":memory:",
		CorsOrigins:     []string{"*"},
		RateLimitReqs:   1000,
		RateLimitWindow: 60,
	}
	config.AppConfig = cfg

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	if err := migrator.RunPending(db); err != nil {
		t.Fatalf("Failed to run migrations on test db: %v", err)
	}

	r := router.New(cfg)
	routes.RegisterRoutes(r)
	return r
}

func TestAuthAndUserEndpoints(t *testing.T) {
	r := setupTestApp(t)

	// 1. Test Register
	regPayload := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(regPayload)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created on register, got %d. Body: %s", w.Code, w.Body.String())
	}

	var regResponse struct {
		Success bool `json:"success"`
		Item    struct {
			Token string `json:"token"`
		} `json:"item"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &regResponse); err != nil || regResponse.Item.Token == "" {
		t.Fatalf("Failed to parse token from register response: %v", err)
	}
	token := regResponse.Item.Token

	// 2. Test Login
	loginPayload := map[string]string{
		"email":    "john@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on login, got %d", w.Code)
	}

	// 3. Test Me with Bearer Token
	req = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on /auth/me, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 4. Test User CRUD Index
	req = httptest.NewRequest(http.MethodGet, "/users?page=1&per_page=10", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on /users, got %d", w.Code)
	}
}
