package models

import (
	"time"

	"github.com/wibiesana/padi-core/auth"
	"github.com/wibiesana/padi-core/model"
	"github.com/wibiesana/padi-core/query"
)

type User struct {
	ID        uint       `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validate:"required"`
	Email     string     `db:"email" json:"email" validate:"required,email"`
	Password  string     `db:"password" json:"password,omitempty" validate:"required"`
	Role      string     `db:"role" json:"role"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Find finds user by ID
func (User) Find(id interface{}) (*User, error) {
	var u User
	err := query.New("users").Where("id", id).First(&u)
	return &u, err
}

// FindByEmail finds user by email
func (User) FindByEmail(email string) (*User, error) {
	var u User
	err := query.New("users").Where("email", email).First(&u)
	return &u, err
}

// Save saves or updates user
func (u *User) Save() error {
	return model.Save(u)
}

// Delete removes user
func (u *User) Delete() error {
	return model.Delete(u)
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
