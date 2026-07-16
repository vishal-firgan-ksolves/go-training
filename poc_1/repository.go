package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
    Create(ctx context.Context, user User) (bool, error)
	GetAll(ctx context.Context) ([]User, error)
    GetByID(ctx context.Context, id string) (User, bool, error)
    Update(ctx context.Context, id string, user User) (User, bool, error)
    Delete(ctx context.Context, id string) (bool, error)
}

type InMemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]User),
	}
}

func (s *InMemoryUserStore) Create(ctx context.Context, user User) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, idExists := s.users[user.Id]; idExists {
		return false, nil
	}

	for _, existingUser := range s.users {
		if existingUser.Email == user.Email {
			return false, nil
		}
	}
	s.users[user.Id] = user
	return true, nil
}

func (s *InMemoryUserStore) GetAll(ctx context.Context) ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users, nil
}

func (s *InMemoryUserStore) GetByID(ctx context.Context, id string) (User, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	return user, exists, nil
}

func (s *InMemoryUserStore) Update(ctx context.Context, id string, user User) (User, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.users[id]
	if !exists {
		return User{}, false, nil
	}

	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Email != "" {
		existing.Email = user.Email
	}

	s.users[id] = existing
	return existing, true, nil
}

func (s *InMemoryUserStore) Delete(ctx context.Context, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return false, nil
	}

	delete(s.users, id)
	return true, nil
}

type CachedUserRepository struct {
	db    UserRepository
	cache *redis.Client
}

func NewCachedUserRepository(db UserRepository, cache *redis.Client) *CachedUserRepository {
	return &CachedUserRepository{
		db:    db,
		cache: cache,
	}
}

func (c *CachedUserRepository) GetByID(ctx context.Context, id string) (User, bool, error) {
	reqID, _ := ctx.Value(requestIDKey).(string)
	cacheKey := "user:" + id
	cachedData, err := c.cache.Get(ctx, cacheKey).Result()

	log.Printf("DEBUG: CachedUserRepository received ID string: '%s'", id)

	if err == nil {
		var user User
		_ = json.Unmarshal([]byte(cachedData), &user)
		log.Printf("[ID: %s] CACHE HIT: Fetched user %s from Redis", reqID, id)
		return user, true, nil
	}

	user, exists, err := c.db.GetByID(ctx, id)

	log.Printf("DEBUG: DB Map lookup result for '%s' -> exists: %t, err: %v", id, exists, err)

	if err != nil || !exists {
		return user, exists, err
	}

	log.Printf("[ID: %s] CACHE MISS: Fetched user %s from DB", reqID, id)

	userJSON, _ := json.Marshal(user)
	c.cache.Set(ctx, cacheKey, userJSON, 5*time.Minute)
	return user, true, nil
}

func (c *CachedUserRepository) Create(ctx context.Context, user User) (bool,error) {
	return c.db.Create(ctx, user)
}

func (c *CachedUserRepository) GetAll(ctx context.Context) ([]User, error) {
	return c.db.GetAll(ctx)
}

func (c *CachedUserRepository) Update(ctx context.Context, id string, user User) (User, bool, error) {
	return c.db.Update(ctx, id, user)
}

func (c *CachedUserRepository) Delete(ctx context.Context, id string) (bool, error) {
	return c.db.Delete(ctx, id)
}