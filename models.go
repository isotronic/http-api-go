package main

import (
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	AccessToken 	string `json:"token"`
	RefreshToken 	string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken 	string `json:"token"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string 		`json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}