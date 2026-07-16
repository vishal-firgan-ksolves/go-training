package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserStore_CRUD(t *testing.T) {
    store := NewInMemoryUserStore()
    ctx := context.Background()

    testUser := User{
        Id:    "user-1",
        Name:  "Original Name",
        Email: "test@example.com",
    }

    created, err := store.Create(ctx, testUser)
    assert.NoError(t, err)
    assert.True(t, created)

    created, err = store.Create(ctx, testUser)
    assert.NoError(t, err)
    assert.False(t, created)

    fetchedUser, exists, err := store.GetByID(ctx, "user-1")
    assert.NoError(t, err)
    assert.True(t, exists)
    assert.Equal(t, "Original Name", fetchedUser.Name)

    allUsers, err := store.GetAll(ctx)
    assert.NoError(t, err)
    assert.Len(t, allUsers, 1)

    updateData := User{Name: "Updated Name"}
    updatedUser, exists, err := store.Update(ctx, "user-1", updateData)
    assert.NoError(t, err)
    assert.True(t, exists)
    assert.Equal(t, "Updated Name", updatedUser.Name)
    assert.Equal(t, "test@example.com", updatedUser.Email)

    deleted, err := store.Delete(ctx, "user-1")
    assert.NoError(t, err)
    assert.True(t, deleted)

    _, existsAfterDelete, _ := store.GetByID(ctx, "user-1")
    assert.False(t, existsAfterDelete)
}