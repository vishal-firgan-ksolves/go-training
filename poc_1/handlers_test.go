package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user User) (bool, error) {
    args := m.Called(ctx, user)
    return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (User, bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Bool(1), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, user User) (User, bool, error) {
	args := m.Called(ctx, id, user)
	return args.Get(0).(User), args.Bool(1), args.Error(2)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestGetUser_WithMock(t *testing.T) {
	mockRepo := new(MockUserRepository)
	
	dummyRedis := redis.NewClient(&redis.Options{})
	
	mockService := NewUserService(mockRepo, dummyRedis)
	handler := NewUserHandler(mockService)

	mockUser := User{
		Id:        "user-123",
		Name:      "Mocked Alice",
		Email:     "alice.mock@example.com",
		CreatedAt: "2026-07-14T12:00:00Z",
	}

	mockRepo.On("GetByID", mock.Anything, "user-123").Return(mockUser, true, nil)

	r := chi.NewRouter()
	r.Get("/api/users/{id}", handler.getUser)

	req := httptest.NewRequest("GET", "/api/users/user-123", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp UserResponseDTO
	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, mockUser.Id, resp.ID)
	assert.Equal(t, mockUser.Name, resp.Name)
	assert.Equal(t, mockUser.Email, resp.Email)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_ValidationTable(t *testing.T) {
	store := NewInMemoryUserStore()
	dummyRedis := redis.NewClient(&redis.Options{})
	
	service := NewUserService(store, dummyRedis)
	handler := NewUserHandler(service)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedErrSub string 
	}{
		{
			name:           "Successful user creation",
			payload:        `{"name": "John Doe", "email": "john.doe@example.com"}`,
			expectedStatus: http.StatusCreated,
			expectedErrSub: "",
		},
		{
			name:           "Failure - Empty Name",
			payload:        `{"name": "", "email": "john.doe@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrSub: "Name",
		},
		{
			name:           "Failure - Name too short",
			payload:        `{"name": "Jo", "email": "john.doe@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrSub: "Name",
		},
		{
			name:           "Failure - Invalid Email Format",
			payload:        `{"name": "John Doe", "email": "not-an-email"}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrSub: "Email",
		},
		{
			name:           "Failure - Missing Email",
			payload:        `{"name": "John Doe", "email": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrSub: "Email",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(tc.payload))
			rec := httptest.NewRecorder()

			handler.CreateUser(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedErrSub != "" {
				var errResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errResp)
				assert.NoError(t, err)
				assert.Contains(t, errResp.Error, tc.expectedErrSub)
			}
		})
	}
}

func TestUserAPI_IntegrationWithNewServer(t *testing.T) {
	store := NewInMemoryUserStore()
	dummyRedis := redis.NewClient(&redis.Options{})
	
	service := NewUserService(store, dummyRedis)
	handler := NewUserHandler(service)

	r := chi.NewRouter()
	r.Post("/api/users", handler.CreateUser)
	r.Get("/api/users", handler.getAllUsers)

	server := httptest.NewServer(r)
	defer server.Close() 

	postPayload := `{"name": "Bob Builder", "email": "bob@example.com"}`
	resp, err := http.Post(server.URL+"/api/users", "application/json", bytes.NewBufferString(postPayload))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser UserResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, "Bob Builder", createdUser.Name)
	assert.Equal(t, "bob@example.com", createdUser.Email)
	assert.NotEmpty(t, createdUser.ID)

	getResp, err := http.Get(server.URL + "/api/users")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var usersList []UserResponseDTO
	err = json.NewDecoder(getResp.Body).Decode(&usersList)
	getResp.Body.Close()

	assert.NoError(t, err)
	assert.Len(t, usersList, 1)
	assert.Equal(t, createdUser.ID, usersList[0].ID)
	assert.Equal(t, "Bob Builder", usersList[0].Name)
}

func TestUpdateUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    dummyRedis := redis.NewClient(&redis.Options{})
    mockService := NewUserService(mockRepo, dummyRedis)
    handler := NewUserHandler(mockService)

    reqPayload := `{"name": "Bob Updated"}`
    
    existingUser := User{
        Id:    "user-123",
        Name:  "Bob Original",
        Email: "bob@example.com",
    }
    
    mockRepo.On("GetByID", mock.Anything, "user-123").Return(existingUser, true, nil)

    returnedUser := User{
        Id:    "user-123",
        Name:  "Bob Updated",
        Email: "bob@example.com",
    }

    mockRepo.On("Update", mock.Anything, "user-123", mock.Anything).Return(returnedUser, true, nil)

    r := chi.NewRouter()
    r.Use(RequestIDMiddleware)
    r.Put("/api/users/{id}", handler.updateUser)

    req := httptest.NewRequest("PUT", "/api/users/user-123", bytes.NewBufferString(reqPayload))
    rec := httptest.NewRecorder()

    r.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)
    mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	dummyRedis := redis.NewClient(&redis.Options{})
	mockService := NewUserService(mockRepo, dummyRedis)
	handler := NewUserHandler(mockService)

	// Mock the Repository Delete call
	mockRepo.On("Delete", mock.Anything, "user-123").Return(true, nil)

	r := chi.NewRouter()
	r.Use(RequestIDMiddleware)
	r.Delete("/api/users/{id}", handler.deleteUser)

	req := httptest.NewRequest("DELETE", "/api/users/user-123", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockRepo.AssertExpectations(t)
}