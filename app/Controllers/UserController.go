package controllers

import (
	base "padi-template/app/Controllers/Base"
)

type UserController struct {
	base.UserController
}

func NewUserController() *UserController {
	return &UserController{}
}

// Custom actions, overrides, or hooks can be added below
