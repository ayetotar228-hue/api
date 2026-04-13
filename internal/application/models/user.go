package models

import (
	"fmt"
	"strings"
	"time"
)

type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(r.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
