package models

import (
	"time"

	"github.com/wibiesana/padi_go_core/activerecord"
	"github.com/wibiesana/padi_go_core/auth"
)

type User struct {
	ID              uint       `db:"id" json:"id"`
	Name            *string    `db:"name" json:"name,omitempty"`
	Username        *string    `db:"username" json:"username,omitempty"`
	Email           string     `db:"email" json:"email" validate:"required,email"`
	Password        string     `db:"password" json:"password,omitempty" validate:"required"`
	Role            string     `db:"role" json:"role"`
	Status          string     `db:"status" json:"status,omitempty"`
	EmailVerifiedAt *int64     `db:"email_verified_at" json:"email_verified_at,omitempty"`
	RememberToken   *string    `db:"remember_token" json:"remember_token,omitempty"`
	LastLoginAt     *int64     `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedBy       *uint      `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy       *uint      `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt       *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Find finds user by ID
func (User) Find(id interface{}) (*User, error) {
	return activerecord.Find[User](id)
}

// FindOrFail finds user by ID or returns error if not found
func (User) FindOrFail(id interface{}) (*User, error) {
	return activerecord.FindOrFail[User](id)
}

// All retrieves all users
func (User) All(columns ...string) ([]User, error) {
	return activerecord.All[User](columns...)
}

// FindByEmail finds user by email
func (User) FindByEmail(email string) (*User, error) {
	return activerecord.FindBy[User]("email", email)
}

// FindByLogin finds user by username, email, or name
func (User) FindByLogin(login string) (*User, error) {
	var item User
	// Check email first
	err := activerecord.FindBuilder[User]().Where("email", login).First(&item)
	if err == nil && item.ID > 0 {
		return &item, nil
	}

	// Check username
	err = activerecord.FindBuilder[User]().Where("username", login).First(&item)
	if err == nil && item.ID > 0 {
		return &item, nil
	}

	// Check name
	err = activerecord.FindBuilder[User]().Where("name", login).First(&item)
	if err == nil && item.ID > 0 {
		return &item, nil
	}

	return nil, nil
}

// Save saves or updates user
func (u *User) Save() error {
	return activerecord.Save(u)
}

// Delete removes user
func (u *User) Delete() error {
	return activerecord.DeleteModel(u)
}

// SetPassword hashes and sets user password
func (u *User) SetPassword(plainPassword string) error {
	hash, err := auth.HashPassword(plainPassword)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

// VerifyPassword checks if plain password matches hash
func (u *User) VerifyPassword(plainPassword string) bool {
	return auth.CheckPasswordHash(plainPassword, u.Password)
}
