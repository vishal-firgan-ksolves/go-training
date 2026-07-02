package service

import (
	"fmt"
	
	"exercise/internal/repository"
)

func GetUser(id int) (*repository.User, error) {
	
	user, err := repository.FindUser(id)
	
	if err != nil {
		return nil, fmt.Errorf("ServiceLayer(GetUser) failed: %w", err)
	}
	
	return user, nil
}