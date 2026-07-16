package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	// 1. Execution
	user, err := CreateUser(99, "Vishal", 25)

	// 2. Assertions (The Testify Way)
	// Notice how we don't use a single "if" statement!
	assert.NoError(t, err, "Should not return an error for valid inputs")
	assert.NotNil(t, user, "User object should successfully be created")
	
	assert.Equal(t, 99, user.ID, "ID should match exactly")
	assert.Equal(t, "Vishal", user.Name)
	assert.True(t, user.IsActive, "New users must be active by default")
}

func TestCreateUser_UnderageFailure(t *testing.T) {
	// 1. Execution
	user, err := CreateUser(100, "Rahul", 15) // 15 is too young!

	// 2. Assertions
	assert.Error(t, err, "Should return an error for underage user")
	assert.EqualError(t, err, "user must be 18 or older", "Error message must match exactly")
	assert.Nil(t, user, "User object should be nil on failure")
}

func TestGenerateUserJSON(t *testing.T) {
	// 1. Setup
	user := &User{ID: 1, Name: "Amit", Age: 30, IsActive: true}

	// 2. Execution
	jsonStr, err := GenerateUserJSON(user)

	// 3. Assertions
	assert.NoError(t, err)

	// JSONEq is MAGICAL. It ignores spaces and formatting!
	expectedJSON := `{
		"id": 1,
		"name": "Amit",
		"age": 30,
		"is_active": true
	}`
	assert.JSONEq(t, expectedJSON, jsonStr, "JSON structure must match")
}