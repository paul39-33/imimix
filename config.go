package main

import (
	"github.com/google/uuid"
	"github.com/paul39-33/imimix/internal/database"
)

type apiConfig struct {
	dbQueries *database.Queries
	secret    string
}

type UserLogin struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}
