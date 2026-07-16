package main

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func getFastDummyRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        "localhost:9999",
		DialTimeout: time.Millisecond * 1,
		MaxRetries:  -1,
	})
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, getFastDummyRedis())
	ctx := context.Background()

	expectedUser := User{
		Id:    "user-123",
		Name:  "Alice",
		Email: "alice@example.com",
	}

	mockRepo.On("GetByID", ctx, "user-123").Return(expectedUser, true, nil)

	user, exists, err := service.GetUser(ctx, "user-123")

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "Alice", user.Name)

	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo, getFastDummyRedis())
    
    ctx := context.WithValue(context.Background(), requestIDKey, "test-req-777")

    newName := "Alice Updated"
    reqDTO := UpdateUserDto{
        Name: &newName,
    }

    existingUser := User{
        Id:    "user-123",
        Name:  "Alice Original",
        Email: "alice@example.com",
    }
    
    mockRepo.On("GetByID", ctx, "user-123").Return(existingUser, true, nil)

    expectedUpdateData := User{
        Id:    "user-123",
        Name:  "Alice Updated",
        Email: "alice@example.com",
    }
    
    returnedUser := expectedUpdateData

    mockRepo.On("Update", ctx, "user-123", expectedUpdateData).Return(returnedUser, true, nil)

    user, exists, err := service.UpdateUser(ctx, "user-123", reqDTO)

    assert.NoError(t, err)
    assert.True(t, exists)
    assert.Equal(t, "Alice Updated", user.Name)
    assert.Equal(t, "alice@example.com", user.Email)

    mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, getFastDummyRedis())
	
	ctx := context.WithValue(context.Background(), requestIDKey, "test-req-888")

	mockRepo.On("Delete", ctx, "user-123").Return(true, nil)

	deleted, err := service.DeleteUser(ctx, "user-123")

	assert.NoError(t, err)
	assert.True(t, deleted)
	mockRepo.AssertExpectations(t)
}