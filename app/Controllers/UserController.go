package controllers

import (
	"database/sql"
	"net/http"

	models "padi-template/app/Models"
	"github.com/wibiesana/padi-core/query"
	"github.com/wibiesana/padi-core/response"
	"github.com/wibiesana/padi-core/router"
	"github.com/wibiesana/padi-core/validator"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// Index lists records with pagination and search
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	opts := query.ParseOptions(r)
	var records []models.User

	searchColumns := []string{"name", "email", "role"}
	meta, err := query.New("users").Paginate(opts, searchColumns, &records)
	if err != nil {
		response.InternalServerError(w, "Failed to retrieve User list")
		return
	}

	response.Paginated(w, records, meta, "User list retrieved successfully")
}

// All retrieves all users without pagination
func (c *UserController) All(w http.ResponseWriter, r *http.Request) {
	var records []models.User
	err := query.New("users").All(&records)
	if err != nil {
		response.InternalServerError(w, "Failed to retrieve all users")
		return
	}

	response.Success(w, records, "All users retrieved successfully")
}

// Show retrieves a single record by ID
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := router.ParamUint(r, "id")
	if err != nil {
		response.BadRequest(w, "Invalid ID parameter")
		return
	}

	item, err := (models.User{}).Find(id)
	if err != nil || item == nil || item.ID == 0 {
		if err == sql.ErrNoRows {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalServerError(w, "Failed to retrieve User")
		return
	}

	response.Success(w, item, "User retrieved successfully")
}

// Store creates a new record
func (c *UserController) Store(w http.ResponseWriter, r *http.Request) {
	var item models.User
	if errs, err := validator.BindJSON(r, &item); err != nil {
		response.UnprocessableEntity(w, errs, "Validation failed")
		return
	}

	if item.Password != "" {
		if err := item.SetPassword(item.Password); err != nil {
			response.InternalServerError(w, "Failed to hash password")
			return
		}
	}

	if err := item.Save(); err != nil {
		response.InternalServerError(w, "Failed to create User: "+err.Error())
		return
	}

	response.Created(w, item, "User created successfully")
}

// Update updates an existing record
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := router.ParamUint(r, "id")
	if err != nil {
		response.BadRequest(w, "Invalid ID parameter")
		return
	}

	item, err := (models.User{}).Find(id)
	if err != nil || item == nil || item.ID == 0 {
		if err == sql.ErrNoRows {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalServerError(w, "Failed to find User")
		return
	}

	var updateData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if errs, err := validator.BindJSON(r, &updateData); err != nil {
		response.UnprocessableEntity(w, errs, "Validation failed")
		return
	}

	if updateData.Name != "" {
		item.Name = updateData.Name
	}
	if updateData.Email != "" {
		item.Email = updateData.Email
	}
	if updateData.Role != "" {
		item.Role = updateData.Role
	}
	if updateData.Password != "" {
		if err := item.SetPassword(updateData.Password); err != nil {
			response.InternalServerError(w, "Failed to hash password")
			return
		}
	}

	if err := item.Save(); err != nil {
		response.InternalServerError(w, "Failed to update User")
		return
	}

	response.Success(w, item, "User updated successfully")
}

// Destroy deletes a record
func (c *UserController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, err := router.ParamUint(r, "id")
	if err != nil {
		response.BadRequest(w, "Invalid ID parameter")
		return
	}

	item, err := (models.User{}).Find(id)
	if err != nil || item == nil || item.ID == 0 {
		if err == sql.ErrNoRows {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalServerError(w, "Failed to find User")
		return
	}

	if err := item.Delete(); err != nil {
		response.InternalServerError(w, "Failed to delete User")
		return
	}

	response.Success(w, nil, "User deleted successfully")
}
