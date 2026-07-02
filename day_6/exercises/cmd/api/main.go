package main

import (
	"errors"
	"fmt"
	"exercise/internal/customerrors"
	"exercise/internal/repository"
	"exercise/internal/service"
)

// We define a function type that matches service.GetUser
type UserFetcher func(id int) (*repository.User, error)

// processUserRequest now accepts the dependency and returns the status string
func processUserRequest(id int, fetchUser UserFetcher) string {
	user, err := fetchUser(id)

	if err == nil {
		return fmt.Sprintf("Status: 200 OK | Success! Found User '%s' (ID: %d)", user.Name, user.ID)
	}

	var validationErr *customerrors.ValidationError
	var notFoundErr *customerrors.NotFoundError

	if errors.As(err, &validationErr) {
		return fmt.Sprintf("Status: 400 Bad Request | Field '%s': %s", validationErr.Field, validationErr.Reason)
	}

	if errors.As(err, &notFoundErr) {
		return fmt.Sprintf("Status: 404 Not Found | Missing %s ID %d", notFoundErr.Resource, notFoundErr.ID)
	}

	if errors.Is(err, customerrors.ErrDatabaseTimeout) {
		return "Status: 503 Service Unavailable | Database connection timed out."
	}

	return fmt.Sprintf("Status: 500 Internal Server Error | %v", err)
}

func main() {
	// In production, we inject the REAL service.GetUser function
	fmt.Println(processUserRequest(10, service.GetUser))
	fmt.Println(processUserRequest(-5, service.GetUser))
	fmt.Println(processUserRequest(404, service.GetUser))
	fmt.Println(processUserRequest(500, service.GetUser))
}