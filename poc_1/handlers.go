package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_= json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var reqDTO CreateUserDto
    if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
        sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }

    if err := Validate.Struct(reqDTO); err != nil {
        sendJSONError(w, err.Error(), http.StatusBadRequest)
        return
    }

    newUser := User{
        Id:        uuid.New().String(),
        Name:      reqDTO.Name,
        Email:     reqDTO.Email,
        CreatedAt: time.Now().Format(time.RFC3339),
    }

    ctx := r.Context()
    created, err := h.service.CreateUser(ctx, newUser)
    if err != nil {
        sendJSONError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if !created {
        sendJSONError(w, "User with this email already exists", http.StatusConflict) // 409 Conflict
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(UserResponseDTO{
        ID:        newUser.Id,
        Name:      newUser.Name,
        Email:     newUser.Email,
        CreatedAt: newUser.CreatedAt,
    })
}

// /api/users
func (h *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	users, err := h.service.GetAllUsers(ctx)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]UserResponseDTO, 0, len(users))
	for _, u := range users {
		resp = append(resp, UserResponseDTO{
			ID:        u.Id,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// /api/users/{id}
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	user, exists, err := h.service.GetUser(ctx, id)

	if err != nil {
		sendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserResponseDTO{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

// /api/users/{id}
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var reqDTO UpdateUserDto
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(reqDTO); err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	updated, exists, err := h.service.UpdateUser(ctx, id, reqDTO)

	if err != nil {
		sendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserResponseDTO{
		ID:        updated.Id,
		Name:      updated.Name,
		Email:     updated.Email,
		CreatedAt: updated.CreatedAt,
	})
}

// /api/users/{id}
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	deleted, err := h.service.DeleteUser(ctx, id)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleted {
		sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}