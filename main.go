package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//load environment variables from .env file
	godotenv.Load()
	//secret := os.Getenv("SECRET_KEY")
	db_url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Print("db connected successfully")
	defer db.Close()

}
