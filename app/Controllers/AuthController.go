package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	models "padi-template/app/Models"
	"github.com/wibiesana/padi_go_core/activerecord"
	"github.com/wibiesana/padi_go_core/auth"
	"github.com/wibiesana/padi_go_core/middleware"
	"github.com/wibiesana/padi_go_core/response"
	"github.com/wibiesana/padi_go_core/validator"
)

func generateSecureRandomString(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Login      string `json:"login"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

type AuthResponse struct {
	User          models.User `json:"user"`
	Token         string      `json:"token"`
	RememberToken *string     `json:"remember_token,omitempty"`
	ExpiresIn     *int64      `json:"expires_in,omitempty"`
}

// Register creates a new user account and returns JWT token
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if errs, err := validator.BindJSON(r, &req); err != nil {
		response.UnprocessableEntity(w, errs, "Validation failed")
		return
	}

	// Check if email already exists using ActiveRecord FindByEmail
	if existing, err := (models.User{}).FindByEmail(req.Email); err == nil && existing != nil && existing.ID > 0 {
		response.BadRequest(w, "Email is already registered")
		return
	}

	user := models.User{}
	if req.Name != "" {
		name := req.Name
		user.Name = &name
	}
	user.Email = req.Email
	user.Role = req.Role
	if user.Role == "" {
		user.Role = "user"
	}
	user.Status = "active"

	if err := user.SetPassword(req.Password); err != nil {
		response.InternalServerError(w, "Failed to hash password")
		return
	}

	// Save via ActiveRecord Save()
	if err := user.Save(); err != nil {
		response.InternalServerError(w, "Failed to create user account: "+err.Error())
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		response.InternalServerError(w, "User created but token generation failed")
		return
	}

	response.Created(w, AuthResponse{
		User:  user,
		Token: token,
	}, "User registered successfully")
}

// Login authenticates user with email/username & password and returns JWT token
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if errs, err := validator.BindJSON(r, &req); err != nil {
		response.UnprocessableEntity(w, errs, "Validation failed")
		return
	}

	loginIdentifier := req.Login
	if loginIdentifier == "" {
		if req.Email != "" {
			loginIdentifier = req.Email
		} else if req.Username != "" {
			loginIdentifier = req.Username
		}
	}

	if loginIdentifier == "" {
		response.UnprocessableEntity(w, map[string]string{
			"login": "email or username is required",
		}, "Validation failed")
		return
	}

	user, err := (models.User{}).FindByLogin(loginIdentifier)
	if err != nil || user == nil || user.ID == 0 {
		response.Unauthorized(w, "Invalid email/username or password")
		return
	}

	if !user.VerifyPassword(req.Password) {
		response.Unauthorized(w, "Invalid email/username or password")
		return
	}

	// Update last login
	now := time.Now().Unix()
	user.LastLoginAt = &now

	var rememberToken *string
	var expiresIn *int64

	if req.RememberMe {
		ttl := int64(365 * 24 * 3600) // 1 year
		expiresIn = &ttl
		rt := generateSecureRandomString(32)
		rememberToken = &rt
		user.RememberToken = rememberToken
	}

	_ = user.Save()

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		response.InternalServerError(w, "Failed to generate authentication token")
		return
	}

	response.Success(w, AuthResponse{
		User:          *user,
		Token:         token,
		RememberToken: rememberToken,
		ExpiresIn:     expiresIn,
	}, "Login successful")
}

// Refresh generates a new access token using remember_token
func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RememberToken string `json:"remember_token" validate:"required"`
	}
	if errs, err := validator.BindJSON(r, &req); err != nil {
		response.UnprocessableEntity(w, errs, "Validation failed")
		return
	}

	user, err := activerecord.FindBy[models.User]("remember_token", req.RememberToken)
	if err != nil || user == nil || user.ID == 0 {
		response.Unauthorized(w, "Invalid or expired remember token")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		response.InternalServerError(w, "Failed to generate token")
		return
	}

	ttl := int64(365 * 24 * 3600)
	response.Success(w, AuthResponse{
		User:          *user,
		Token:         token,
		RememberToken: user.RememberToken,
		ExpiresIn:     &ttl,
	}, "Token refreshed successfully")
}

// Logout invalidates user session
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	response.Success(w, nil, "Logout successful")
}

// Me returns currently authenticated user profile
func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok || claims == nil {
		response.Unauthorized(w, "Unauthenticated")
		return
	}

	user, err := (models.User{}).Find(claims.UserID)
	if err != nil {
		response.InternalServerError(w, "Failed to retrieve user profile")
		return
	}
	if user == nil || user.ID == 0 {
		response.NotFound(w, "User profile not found")
		return
	}

	response.Success(w, user, "User profile retrieved")
}
