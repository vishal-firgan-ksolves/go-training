package main

import (
	"encoding/json"
	"errors"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
}

func CreateUser(id int, name string, age int) (*User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if age < 18 {
		return nil, errors.New("user must be 18 or older")
	}

	return &User{
		ID:       id,
		Name:     name,
		Age:      age,
		IsActive: true,
	}, nil
}

func GenerateUserJSON(u *User) (string, error) {
	bytes, err := json.Marshal(u)
	return string(bytes), err
}

func main() {
}