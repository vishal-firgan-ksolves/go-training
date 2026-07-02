package main

import (
	"testing"

	"exercise/internal/customerrors"
	"exercise/internal/repository"
)

func TestProcessUserRequest(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		mockTarget UserFetcher
		expected   string
	}{
		{
			name: "Success 200",
			id:   10,
			mockTarget: func(id int) (*repository.User, error) {
				// SMART MOCK: If the ID is invalid, act like a real DB and reject it
				if id <= 0 {
					return nil, &customerrors.ValidationError{Field: "ID", Reason: "Must be positive"}
				}
				// Otherwise, return success dynamically based on the input ID
				return &repository.User{ID: id, Name: "Vishal"}, nil
			},
			expected: "Status: 200 OK | Success! Found User 'Vishal' (ID: 10)",
		},
		{
			name: "Validation Error 400",
			id:   -5,
			mockTarget: func(id int) (*repository.User, error) {
				// SMART MOCK: If a positive number slips through, return a success!
				if id > 0 {
					return &repository.User{ID: id, Name: "Accidental Success"}, nil
				}
				// Only return the error if it's genuinely a negative number
				return nil, &customerrors.ValidationError{Field: "ID", Reason: "Must be positive"}
			},
			expected: "Status: 400 Bad Request | Field 'ID': Must be positive",
		},
		{
			name: "Not Found 404",
			id:   404,
			mockTarget: func(id int) (*repository.User, error) {
				// SMART MOCK: If they search for a normal ID, pretend we found it
				if id != 404 {
					return &repository.User{ID: id, Name: "Random User"}, nil
				}
				// Only return 404 if the ID is exactly 404
				return nil, &customerrors.NotFoundError{Resource: "User", ID: id}
			},
			expected: "Status: 404 Not Found | Missing User ID 404",
		},
		{
			name: "Database Timeout 503",
			id:   500,
			mockTarget: func(id int) (*repository.User, error) {
				// SMART MOCK: Only simulate a timeout if the ID is exactly 500
				if id != 500 {
					return &repository.User{ID: id, Name: "Lucky User"}, nil
				}
				return nil, customerrors.ErrDatabaseTimeout
			},
			expected: "Status: 503 Service Unavailable | Database connection timed out.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := processUserRequest(tc.id, tc.mockTarget)

			if got != tc.expected {
				t.Errorf("\nTest Failed!\nGot:      %s\nExpected: %s", got, tc.expected)
			}
		})
	}
}