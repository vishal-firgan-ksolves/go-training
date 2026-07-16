package main

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type CreateUserDto struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserDto struct {
	Name  *string `json:"name" validate:"omitempty,min=3,max=50"`
	Email *string `json:"email" validate:"omitempty,email"`
}

type UserResponseDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}