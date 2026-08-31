package models

import (
	base "padi-template/app/Models/Base"
	"github.com/wibiesana/padi-core/auth"
	"github.com/wibiesana/padi-core/model"
	"github.com/wibiesana/padi-core/query"
)

type User struct {
	base.User
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
