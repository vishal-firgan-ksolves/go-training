package main

import (
	"errors"
	"time"
	"fmt"
)

type DatabaseError struct{
	Query string
	Timestamp time.Time
	Err error
}

func (dbErr *DatabaseError) Error() string{
	return fmt.Sprintf("Database Error at %s | Query [%s] Details %v",
	dbErr.Timestamp.Format(time.RFC3339),
	dbErr.Query,
	dbErr.Err)
}

func ExecuteQuery(query string) error {
	// simulatedDBFailure := fmt.Errorf("connection refused")
	customError := errors.New("The database level error")

	// Instantiate and return the custom error struct
	return &DatabaseError{
		Query: query,
		Timestamp: time.Now(),
		Err: customError,
	}
}

func main() {
	err := ExecuteQuery("SELECT * FROM users")
	
	if err != nil {
		// The fmt package automatically calls our custom Error() method
		fmt.Println(err)
	}
}