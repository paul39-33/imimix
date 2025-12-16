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

type User struct {
	ID	   uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Pass string    `json:"pass"`
}

type MimixObj struct {
	ID      uuid.UUID `json:"id"`
	Obj   	string    `json:"obj"`
	ObjType string    `json:"obj_type"`
	PromoteDate time.Time `json:"promote_date"`
	ObjVer  string    `json:"obj_ver"`
	Lib 	string   `json:"lib"`
	MimixStatus string    `json:"mimix_status"`
}

type MimixLib struct {
	ID      uuid.UUID `json:"id"`
	Lib	 string    `json:"lib"`
}


