package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	dbQueries := database.New(db)

	//apiCfg
	apiCfg := apiConfig{
		dbQueries: dbQueries,
		secret:    secret,
	}

	//create Gin router
	r := gin.Default()

	//set frontend directory by url /app/... from .
	r.StaticFS("/app", gin.Dir(".", false))

	api := r.Group("/api")
	{
		api.POST("/login", apiCfg.handleUserLogin)
		api.POST("/add_mimix_lib", apiCfg.handleAddLib)
		api.POST("/add_mimix_obj", apiCfg.handleAddObj)
		api.DELETE("/delete_mimix_lib/:libid", apiCfg.handleDeleteLib)
		api.DELETE("/delete_mimix_obj/:objid", apiCfg.handleDeleteObj)
	}

	//start server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
