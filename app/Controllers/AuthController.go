package controllers

import (
	"database/sql"
	"net/http"

	models "padi-template/app/Models"
	"github.com/wibiesana/padi-core/auth"
	"github.com/wibiesana/padi-core/middleware"
	"github.com/wibiesana/padi-core/response"
	"github.com/wibiesana/padi-core/validator"
)

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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
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
	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	if user.Role == "" {
		user.Role = "user"
	}

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

// Login authenticates user with email & password and returns JWT token
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if errs, err := validator.BindJSON(r, &req); err != nil {
		response.UnprocessableEntity(w, errs, "Validation failed")
		return
	}

	user, err := (models.User{}).FindByEmail(req.Email)
	if err != nil || user == nil || user.ID == 0 {
		response.Unauthorized(w, "Invalid email or password")
		return
	}

	if !user.VerifyPassword(req.Password) {
		response.Unauthorized(w, "Invalid email or password")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		response.InternalServerError(w, "Failed to generate authentication token")
		return
	}

	response.Success(w, AuthResponse{
		User:  *user,
		Token: token,
	}, "Login successful")
}

// Me returns currently authenticated user profile
func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.JWTClaims)
	if !ok || claims == nil {
		response.Unauthorized(w, "Unauthenticated")
		return
	}

	user, err := (models.User{}).Find(claims.UserID)
	if err != nil || user == nil || user.ID == 0 {
		if err == sql.ErrNoRows {
			response.NotFound(w, "User profile not found")
			return
		}
		response.InternalServerError(w, "Failed to retrieve user profile")
		return
	}

	response.Success(w, user, "User profile retrieved")
}
