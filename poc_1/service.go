package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	repo  UserRepository
	cache *redis.Client
}

func NewUserService(repo UserRepository, cache *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (User, bool, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, user User) (bool, error) {
    created, err := s.repo.Create(ctx, user)
    if err != nil {
        return false, err
    }
    
    return created, nil
}
func (s *UserService) UpdateUser(ctx context.Context, id string, reqDTO UpdateUserDto) (User, bool, error) {
    reqID, _ := ctx.Value(requestIDKey).(string)

    existingUser, exists, err := s.repo.GetByID(ctx, id)
    if err != nil || !exists {
        return User{}, false, err
    }

    if reqDTO.Name != nil {
        existingUser.Name = *reqDTO.Name
    }

    if reqDTO.Email != nil {
        existingUser.Email = *reqDTO.Email
    }

    updatedUser, exists, err := s.repo.Update(ctx, id, existingUser)
    if err != nil || !exists {
        return updatedUser, exists, err
    }

    cacheKey := "user:" + id
    s.cache.Del(ctx, cacheKey)
    
    log.Printf("[ID: %s] CACHE INVALIDATED: Purged user %s from Redis", reqID, id)

    return updatedUser, true, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) (bool, error) {
	reqID, _ := ctx.Value(requestIDKey).(string)

	deleted, err := s.repo.Delete(ctx, id)
	if err != nil || !deleted {
		return deleted, err
	}

	cacheKey := "user:" + id

	s.cache.Del(ctx, cacheKey)

	log.Printf("[ID: %s] CACHE INVALIDATED: Purged user %s from Redis", reqID, id)

	return true, nil
}