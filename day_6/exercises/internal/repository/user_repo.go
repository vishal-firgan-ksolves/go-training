package repository

import (
	"exercise/internal/customerrors"
)

type User struct {
	ID   int
	Name string
}

func FindUser(id int) (*User, error) {
	if id < 0 {
		return nil, &customerrors.ValidationError{Field: "id", Reason: "ID cannot be negative"}
	}
	if id == 404 {
		return nil, &customerrors.NotFoundError{Resource: "User", ID: id}
	}
	if id == 500 {
		return nil, customerrors.ErrDatabaseTimeout
	}
	
	mockUser := &User{
		ID:   id,
		Name: "Vishal Golang Developer",
	}
	return mockUser, nil
}